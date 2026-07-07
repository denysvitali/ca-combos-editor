# NV ITEM 00028874 (`RFNV_LTE_CA_BW_CLASS_COMBO_I`) Format

This document describes the on-disk format of Qualcomm NV item 00028874, used to store the LTE carrier aggregation (CA) bandwidth class combinations supported by a device. The description is based on the reverse engineering done in this project and covers the zlib wrapper, the 4-byte uncompressed header, the supported entry types, the band record layouts, and the sorting/ordering semantics observed in the field.

## File Layout

A 00028874 file is a zlib-compressed binary payload. The overall structure is:

```
┌──────────────────────────────────────────────────────────────┐
│  zlib compressed data                                          │
│  ┌──────────────────────────────────────────────────────┐    │
│  │  4-byte header (see below)                            │    │
│  ├──────────────────────────────────────────────────────┤    │
│  │  Entry 0: 2-byte type + band records                  │    │
│  │  Entry 1: 2-byte type + band records                  │    │
│  │  ...                                                  │    │
│  └──────────────────────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────┘
```

### 1. Zlib Compression Layer

The 00028874 file itself is a raw zlib (RFC 1950) stream. It is **not** a gzip file and does not contain a gzip header or filename. The only data required is the zlib compressed payload.

To decompress, run the raw file through a zlib decompressor. For example, using `qpdf` or `zlib-flate`:

```bash
zlib-flate --uncompress < 00028874 > extracted.bin
```

To compress a new payload back to the 00028874 format, use level 6 (the level observed in the original files):

```bash
zlib-flate --compress=6 < extracted.bin > 00028874
```

The uncompressed payload is what this tool reads and writes.

### 2. Uncompressed Payload Header

The first four bytes of the decompressed payload are the header:

| Offset | Size | Description |
|--------|------|-------------|
| 0x00   | 1    | Always `0x00` |
| 0x01   | 1    | Always `0x00` |
| 0x02   | 2    | Little-endian unsigned 16-bit entry count |

The two leading zero bytes are constants; the parser expects them to be exactly `00 00`. The remaining two bytes are the number of entries that follow, stored as a little-endian `uint16`.

Example header from a real device file:

```
00 00 c1 00  89 00 07 00 ...
```

`c1 00` = `0x00C1` = 193 entries, followed by the first entry.

### 3. Entry Type Table

After the header, each entry begins with a 2-byte little-endian type discriminator. The high byte is always `0x00`. The known entry types are:

| Value | Name | Direction | Extra per band |
|-------|------|-----------|----------------|
| `137` | `EntryTypeDownlinkNoMIMO` | Downlink | none |
| `138` | `EntryTypeUplinkNoMIMO` | Uplink | none |
| `201` | `EntryTypeDownlinkMIMO` | Downlink | 1-byte MIMO layer count |
| `202` | `EntryTypeUplinkMIMO` | Uplink | 1-byte MIMO layer count |
| `333` | `EntryTypeDownlinkAntennas` | Downlink | 8-byte antenna list |
| `334` | `EntryTypeUplinkAntennas` | Uplink | 8-byte antenna list |

A single 00028874 file normally uses a consistent format for every entry (e.g., 137/138, 201/202, or 333/334). The 333/334 antenna-aware variant is parsed for inspection and can be written with writer mode `333` (`COMBOWRITER_333_334`).

### 4. Band Record Layouts

Every entry contains exactly **6 band slots**, regardless of how many bands are actually used. Unused slots are filled with zero bytes. The size of a slot depends on the entry type.

#### 4.1 137 / 138 (No MIMO)

Each slot is **3 bytes**:

| Field | Size | Description |
|-------|------|-------------|
| Band | 2 bytes | Little-endian unsigned 16-bit band number (e.g., 1, 3, 7, 41) |
| Class | 1 byte | Bandwidth class (A=1, B=2, C=3, D=4, E=5, F=6, G=7, H=8, I=9) |

Slot size: 3 bytes  
Entry payload size: 2 (type) + 6 × 3 = **20 bytes**

#### 4.2 201 / 202 (With MIMO)

Each slot is **4 bytes**:

| Field | Size | Description |
|-------|------|-------------|
| Band | 2 bytes | Little-endian unsigned 16-bit band number |
| Class | 1 byte | Bandwidth class (A=1, B=2, C=3, D=4, E=5, F=6, G=7, H=8, I=9) |
| MIMO | 1 byte | MIMO layer count (commonly 1, 2, or 4) |

Slot size: 4 bytes  
Entry payload size: 2 (type) + 6 × 4 = **26 bytes**

#### 4.3 333 / 334 (With Antenna Lists)

Each slot is **11 bytes**:

| Field | Size | Description |
|-------|------|-------------|
| Band | 2 bytes | Little-endian unsigned 16-bit band number |
| Class | 1 byte | Bandwidth class (A=1, B=2, C=3, D=4, E=5, F=6, G=7, H=8, I=9) |
| Antennas | 8 bytes | Antenna index list (see section 7) |

Slot size: 11 bytes  
Entry payload size: 2 (type) + 6 × 11 = **68 bytes**

### 5. Valid Band Filtering

When reading slots, a slot is considered populated only if:

- The band number is between 1 and 255 inclusive, and
- The class value is between 1 and 9 inclusive.

Any slot that fails either test is ignored. This is how empty slots are skipped in practice. A valid slot with band `0` or class `0` is still treated as empty.

### 6. Sorting Behavior

The editor maintains a specific order when writing entries to the binary format and when displaying them as text.

#### 6.1 Downlink Bands

Inside a downlink entry, bands are sorted **before serialization** by the writer so that:

- The band number is in descending order (highest band first).
- If two bands share the same number, the class is in descending order (largest class first).

For example, the human-readable text `1A4-3A2` becomes `3A2-1A4` when serialized, because 3 > 1. The band `2A2-46E2-48C2` is reordered to `48C2-46E2-2A2` (48, 46, 2).

The text rendering also sorts using the same descending rule, so the output order matches the binary order.

#### 6.2 Uplink Entries

When a single downlink has multiple uplink entries (e.g., from a `downlink.txt`/`uplink.txt` pair), the uplink entries are sorted by their **first band number in ascending order** before being written.

Inside each uplink entry, bands are stored in the order they were given and are also rendered in text using the same descending band/class rule as downlink entries.

### 7. Antenna List Semantics

The 333/334 format stores an antenna port list for each populated band. The list is serialized as **8 raw bytes** immediately after the class byte. Each byte is interpreted as follows:

- A **non-zero** byte is an active antenna index. The value is treated as a **1-based antenna port number** (e.g., `0x01` = antenna 1, `0x02` = antenna 2).
- A **zero** byte is padding and is ignored.
- The eight bytes are read in order; the resulting `[]Antenna` slice preserves that order.
- The writer emits up to eight antennas from `Band.Antennas`, placing them in the first bytes of the 8-byte field and zero-padding the remainder. Extra antennas beyond the eighth are dropped.

Example: if the 8 antenna bytes are `01 02 03 00 00 00 00 00`, the parsed antenna list is `[1, 2, 3]`.

There is no additional structure such as bitmasks or counts; each populated byte is treated as a single 1-based antenna index.

### 8. Reverse-Engineering Observations from Fixtures

The files under `test/resources/` show the following consistent patterns:

- **Entry-type families are not mixed.** Each fixture uses exactly one pair of entry types: either 137/138, 201/202, or (in the 010 Editor template) 333/334. A real 00028874 file never interleaves families.
- **DL-then-UL grouping.** Within a family, every downlink entry is followed immediately by one or more uplink entries. Some downlinks have a single uplink; others have several (e.g., the 2019-11-26 fixtures show `[201, 202, 202]` blocks).
- **Zlib wrapper.** Every compressed fixture begins with the two-byte zlib header `78 9c`. This means CMF = `0x78` (deflate, 32 KiB window) and FLG = `0x9c` (FLEVEL = `10`, i.e., default compression, no preset dictionary). The FLEVEL value matches the zlib level 6 used by the compressor helper and the `compress.sh` script.
- **Original payloads are not pre-sorted.** Some device-extracted downlink entries store bands in ascending order (e.g., `1A-5A` or `3A-7C`). The writer normalizes every downlink entry to descending band number, then descending class, before serialization.
- **Stable header.** All decompressed payloads start with two zero bytes followed by a little-endian `uint16` entry count.

### 9. Human-Readable Text Format

The tool supports two input formats: a single `bands.txt` file and a split `downlink.txt`/`uplink.txt` pair.

#### 9.1 Single bands.txt format

One combo per line. A line consists of a downlink combo followed by an optional uplink combo. The syntax for a single component is:

```
<BAND><CLASS><MIMO?><ULCLASS?>
```

Examples:

- `3A2` → band 3, class A, MIMO 2 (downlink only)
- `1A4A` → band 1, class A, MIMO 4 in DL, plus band 1, class A (MIMO 1) in UL
- `3A2A-1A4` → DL: 3A2, 1A4; UL: 1A

The optional uplink class is the second class letter after the MIMO count. If the MIMO count is omitted, it defaults to 1.

#### 9.2 Split downlink.txt / uplink.txt format

`downlink.txt` contains one downlink combo per line. `uplink.txt` contains a comma-separated list of uplink combos per line, with one line per downlink line. Uplink lists are sorted by the first band number before being written to the binary file.

### 10. Safety Notes

This file is read by the modem at boot. An invalid 00028874 can cause the modem to reject the NV item, fail to register, or bootloop. Always keep a backup of the original file and monitor `dmesg` for modem failures when testing a modified file.

### 11. Native `compress` / `decompress` Commands

The tool can manage the zlib wrapper directly, so the `zlib-flate` workflow is optional:

- `ca-combos-editor decompress <00028874> <extracted.bin>` — read a raw zlib stream and write the uncompressed payload.
- `ca-combos-editor compress <extracted.bin> <00028874>` — compress an uncompressed payload with zlib level 6 and write the raw 00028874 file.

Both commands use the same level-6 settings observed in original device files and are exact replacements for the `compress.sh` / `uncompress.sh` scripts.
