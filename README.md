# Didban – Music Fetching Library for Sorn

**Didban** (Persian: دیدبان – *watchman*) is a Go library for downloading and tagging high-quality music tracks, built for integration into the **Sorn Music System**.

It uses **Deezer track IDs** as the authoritative source of metadata and audio location, pulling content either from **Deezer** (direct encrypted audio) or optionally falling back to **YouTube** if needed.

---

## ✨ What It Does

- 🎵 Downloads songs using a **Deezer track ID**
- 📥 Fetches from:
  - **Deezer**: via official APIs (FLAC / MP3_320 / MP3_128)
  - **YouTube**: fallback (not implemented in this repo directly)
- 🔐 Decrypts Deezer streams using **Blowfish** decryption
- 🏷️ Tags the downloaded audio with complete metadata
- 🧠 Designed for internal use in **Sorn’s backend systems**

---

## 🔧 Features

- Fetches encrypted media URLs from Deezer using a license token
- Decrypts every third audio chunk using a Blowfish key derived from the track ID
- Applies ID3 tags including artist, title, album, and cover art
- Outputs `.mp3` files with clean filenames
- Handles login/authentication via ARL cookie

---

## 📦 Installation

```bash
go get github.com/nawfay/didban
```