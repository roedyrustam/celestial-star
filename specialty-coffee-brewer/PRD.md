# Product Requirements Document (PRD) - Celestial Brew

## 1. Project Vision
**Celestial Brew** adalah aplikasi evaluasi seduhan kopi specialty yang menggabungkan presisi ilmiah dengan estetika visual yang premium. Tujuannya adalah membantu *home brewer* dan profesional mencapai "Perfect Brew" melalui analisis parameter yang mendalam.

## 2. Target Audience
- **Home Brewers**: Penggemar kopi yang ingin meningkatkan kualitas seduhan manual mereka.
- **Coffee Enthusiasts**: Pengguna yang menyukai data dan statistik tentang kopi mereka.
- **Baristas**: Alat bantu cepat untuk kalibrasi rasa di bar.

## 3. Core Features
### 3.1. Brewing Scoring Engine (Go)
- Algoritma penilaian berdasarkan standar SCA (Specialty Coffee Association).
- Mendukung 5 metode: V60/Pourover, Espresso, AeroPress, French Press, Cold Brew.
- Validasi dosis kopi, yield air, suhu, dan waktu seduh.

### 3.2. Flavor Profile Prediction
- Prediksi karakter rasa berdasarkan data input (Acidity, Sweetness, Body, Bitterness).
- Feedback tekstual yang memberikan saran perbaikan teknis (misal: "haluskan gilingan").

### 3.3. Premium Interactive UI
- Tema "Celestial" yang dinamis dengan animasi bintang.
- Brew Timer terintegrasi.
- Kalkulator rasio air:kopi otomatis.

### 3.4. Persistence & Sharing
- Penyimpanan riwayat (History) hingga 10 seduhan terakhir secara lokal.
- Fitur "Copy Recipe" untuk berbagi resep ke media sosial dengan format yang rapi.

## 4. Technical Constraints
- Backend: Go (Standard Library).
- Frontend: Vanilla HTML, CSS, JavaScript (Tanpa framework luar).
- Storage: Browser LocalStorage (Offline-First).

## 5. Future Roadmap
- Integrasi Cloud Persistence (Supabase/Firebase).
- Fitur ekspor laporan seduhan dalam format PDF/Image.
- Komunitas resep global (User-generated recipes).
