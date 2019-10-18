#!/bin/bash
# Compress to 00028874 format
zlib-flate -compress=6 < $1 > 00028874
