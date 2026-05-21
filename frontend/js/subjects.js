const subjects = {
  TYT: {
    name: "TYT",
    topics: {
      "Türkçe": ["Sözcükte Anlam", "Cümlede Anlam", "Paragraf", "Ses Bilgisi", "Yazım Kuralları", "Noktalama İşaretleri", "Sözcük Türleri", "Fiiller", "Cümlenin Ögeleri", "Cümle Çeşitleri", "Anlatım Bozukluğu"],
      "Matematik": ["Temel Kavramlar", "Sayı Basamakları", "Bölme-Bölünebilme", "EBOB-EKOK", "Rasyonel Sayılar", "Basit Eşitsizlikler", "Mutlak Değer", "Üslü Sayılar", "Köklü Sayılar", "Çarpanlara Ayırma", "Oran-Orantı", "Denklem Çözme", "Problemler", "Kümeler", "Fonksiyonlar", "Permütasyon", "Kombinasyon", "Olasılık", "İstatistik"],
      "Fizik": ["Fizik Bilimine Giriş", "Madde ve Özellikleri", "Hareket ve Kuvvet", "Enerji", "Isı ve Sıcaklık", "Elektrik", "Manyetizma", "Dalgalar", "Optik"],
      "Kimya": ["Kimya Bilimi", "Atom ve Periyodik Sistem", "Kimyasal Türler Arası Etkileşimler", "Maddenin Halleri", "Doğa ve Kimya", "Asitler ve Bazlar", "Kimya Her Yerde"],
      "Biyoloji": ["Yaşam Bilimi Biyoloji", "Canlıların Yapısında Bulunan Organik Bileşikler", "Hücre", "Canlıların Çeşitliliği ve Sınıflandırılması", "Hücre Bölünmeleri", "Kalıtım", "Ekoloji"],
      "Tarih": ["Tarih ve Zaman", "İnsanlığın İlk Dönemleri", "Orta Çağ'da Dünya", "İlk ve Orta Çağ'da Türkler", "İslamiyet'in Doğuşu", "Türklerin İslamiyet'i Kabulü", "Beylikten Devlete", "Dünya Gücü Osmanlı", "Değişen Dünya ve Avrupa", "Osmanlı Kültür ve Medeniyeti"],
      "Coğrafya": ["Doğa ve İnsan", "Dünya'nın Şekli ve Hareketleri", "Coğrafi Koordinat Sistemi", "Harita Bilgisi", "Atmosfer ve İklim", "Yeryüzündeki İklim Tipleri", "Topoğrafya ve Kayaçlar", "İç Kuvvetler", "Dış Kuvvetler", "Nüfus", "Göç", "Türkiye'de Yerleşme"],
      "Felsefe": ["Felsefenin Alanı", "Bilgi Felsefesi", "Varlık Felsefesi", "Ahlak Felsefesi", "Sanat Felsefesi", "Din Felsefesi", "Siyaset Felsefesi", "Bilim Felsefesi"],
      "Din Kültürü": ["İnsan ve Din", "Allah İnancı", "İbadet", "Kur'an'da Bazı Kavramlar", "İslam ve Toplum", "Haklar ve Özgürlükler"]
    }
  },
  AYT: {
    name: "AYT",
    topics: {
      "Matematik": ["Trigonometri", "Logaritma", "Diziler", "Limit", "Türev", "İntegral", "Matris-Determinant", "Karmaşık Sayılar", "Polinomlar", "Çember ve Daire", "Katı Cisimler", "Analitik Geometri", "Olasılık ve İstatistik"],
      "Türk Dili ve Edebiyatı": ["Türk Edebiyatının Dönemleri", "İslamiyet Öncesi Türk Edebiyatı", "İslami Dönem Türk Edebiyatı", "Halk Edebiyatı", "Divan Edebiyatı", "Tanzimat Edebiyatı", "Servetifünun Edebiyatı", "Fecriati Edebiyatı", "Milli Edebiyat", "Cumhuriyet Dönemi Edebiyatı", "Günümüz Türk Edebiyatı"],
      "Fizik": ["Vektörler", "Bağıl Hareket", "Dinamik", "İş-Güç-Enerji", "Atışlar", "Dönme Hareketi", "Elektrik Alan", "Manyetik Alan", "İndüksiyon", "Alternatif Akım", "Dalga Mekaniği", "Atom Fiziği", "Radyoaktivite"],
      "Kimya": ["Modern Atom Teorisi", "Gazlar", "Sıvı Çözeltiler", "Kimyasal Tepkimelerde Enerji", "Tepkime Hızları", "Kimyasal Denge", "Asit-Baz Dengesi", "Çözünürlük Dengesi", "Elektrokimya"],
      "Biyoloji": ["Sinir Sistemi", "Endokrin Sistem", "Duyu Sistemleri", "Destek ve Hareket Sistemi", "Dolaşım Sistemi", "Solunum Sistemi", "Boşaltım Sistemi", "Üreme Sistemi", "Bağışıklık Sistemi", "Bitki Biyolojisi", "Canlılar ve Çevre", "DNA ve Genetik", "Evrim"],
      "Tarih": ["Tarih Bilimi", "İlk Çağ Uygarlıkları", "Türk Tarihi", "Osmanlı Tarihi", "Avrupa Tarihi", "Yakın Çağ Dünya Tarihi", "Türkiye Cumhuriyeti Tarihi", "Soğuk Savaş Dönemi", "Küreselleşen Dünya"],
      "Coğrafya": ["Doğal Sistemler", "Beşeri Sistemler", "Bölgeler ve Ülkeler", "Küresel Ortam", "Türkiye'nin Coğrafi Konumu", "Türkiye'de İklim", "Türkiye'de Yer Şekilleri", "Türkiye'de Nüfus", "Türkiye'de Tarım", "Türkiye'de Sanayi"],
      "Felsefe": ["Psikoloji Bilimi", "Davranış ve Süreçleri", "Öğrenme", "Bellek", "Düşünme ve Dil", "Sosyoloji Bilimi", "Toplumsal Yapı", "Toplumsal Değişme", "Kültür", "Toplumsal Kurumlar"],
      "Din Kültürü": ["Vahiy ve Akıl", "İnanç ve İbadet", "İslam Ahlakı", "Kur'an'da İnsan", "Hak ve Sorumluluk", "İslam ve Bilim"]
    }
  }
};
