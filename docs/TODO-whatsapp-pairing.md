# TODO: WhatsApp Pairing Issue

## Status: PENDING - Akan diperbaiki setelah fungsi lain selesai

## Masalah
Saat mencoba pairing WhatsApp, device menampilkan error: **"tidak dapat menautkan perangkat"**

## Gejala
- QR code berhasil di-generate
- User bisa scan QR code dengan WhatsApp di HP
- Tapi WhatsApp menolak pairing dengan pesan error

## Kemungkinan Penyebab
1. **Limit Linked Devices**: WhatsApp hanya mengizinkan maksimal 4 linked devices. Cek apakah sudah ada 4 devices terpasang
2. **Whatsmeow Configuration**: Mungkin perlu konfigurasi khusus untuk client ID atau device info
3. **WhatsApp API Restrictions**: Ada pembatasan dari WhatsApp untuk certain cases
4. **Session/Store Issue**: Kemungkinan issue dengan database session storage

## Perbaikan yang Sudah Dilakukan
✅ Remove manual phone number input
✅ Auto-detect phone number from QR scan
✅ Keep WhatsApp connection alive in background goroutine
✅ Automatic phone number save to database

## Next Steps untuk Debugging
1. Cek berapa linked devices yang sudah ada di WhatsApp
2. Review whatsmeow documentation untuk proper client initialization
3. Coba dengan fresh WhatsApp session (unlink devices yang lain)
4. Debug dengan whatsmeow verbose logging
5. Cek apakah perlu custom Client ID/Device Name

## Workaround Sementara
Untuk sementara, testing WhatsApp bot bisa menggunakan:
- Test dengan WhatsApp account yang belum memiliki linked devices
- Atau unlink devices lain dulu untuk free up slot

## Files Terkait
- `internal/whatsapp/client.go` - WhatsApp client implementation
- `internal/handler/whatsapp_handler.go` - Pairing handler
- `templates/admin/whatsapp_new.html` - UI for pairing

## Priority
Medium - Fungsi lain akan diselesaikan dulu, baru kembali fix pairing issue ini

## Created
2025-11-16

## Notes
Error terjadi pada saat scan QR code, bukan pada generate QR code.
Connection sudah stay alive di background, tapi WhatsApp server reject pairing.
