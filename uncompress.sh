#!/bin/bash
# Compress to 00028874 format
zlib-flate -uncompress < $1 > extracted
