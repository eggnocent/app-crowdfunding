-- +migrate Up
INSERT INTO campaigns (user_id, name, short_description, description, goal_amount, current_amount, perks, backer_count, slug, created_at, updated_at)
VALUES
(1, 'Bantu Pendidikan Anak Indonesia', 'Membantu pendidikan anak di daerah terpencil', 'Penggalangan dana untuk memberikan akses pendidikan yang lebih baik bagi anak-anak di daerah terpencil di Indonesia.', 500000000, 250000000, 'Buku gratis, Seragam sekolah, Akses internet', 500, 'bantu-pendidikan-anak-indonesia', NOW(), NOW()),
(2, 'Peduli Lingkungan', 'Menanam 1000 pohon untuk lingkungan lebih baik', 'Inisiatif untuk menanam 1000 pohon di berbagai wilayah Indonesia untuk mengurangi dampak perubahan iklim.', 200000000, 80000000, 'Tanam pohon, Sertifikat penghargaan, Merchandise', 300, 'peduli-lingkungan', NOW(), NOW()),
(2, 'Bantuan Bencana Alam', 'Bantu korban bencana alam di Indonesia', 'Penggalangan dana untuk membantu korban bencana alam seperti gempa, banjir, dan letusan gunung berapi di berbagai daerah di Indonesia.', 1000000000, 750000000, 'Paket sembako, Pakaian layak pakai, Obat-obatan', 1000, 'bantuan-bencana-alam', NOW(), NOW()),
(1, 'Pembangunan Sarana Air Bersih', 'Menyediakan air bersih untuk daerah kekurangan air', 'Proyek penggalangan dana untuk pembangunan fasilitas air bersih di desa-desa yang mengalami kekurangan air.', 300000000, 150000000, 'Akses air bersih, Seminar kesehatan, Merchandise', 400, 'pembangunan-sarana-air-bersih', NOW(), NOW()),
(2, 'Rumah Singgah Anak Yatim', 'Membangun rumah singgah untuk anak yatim', 'Penggalangan dana untuk membangun rumah singgah bagi anak yatim di Jakarta agar mereka memiliki tempat tinggal yang aman.', 700000000, 450000000, 'Tempat tinggal, Pendidikan, Makanan bergizi', 700, 'rumah-singgah-anak-yatim', NOW(), NOW());

-- +migrate Down
DELETE FROM campaigns WHERE slug IN ('bantu-pendidikan-anak-indonesia', 'peduli-lingkungan', 'bantuan-bencana-alam', 'pembangunan-sarana-air-bersih', 'rumah-singgah-anak-yatim');
