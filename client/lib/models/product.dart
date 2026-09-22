class Product {
  final String id;
  final String name;
  final double price;
  final String category;

  Product({
    required this.id,
    required this.name,
    required this.price,
    required this.category,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      // 1. Folosim 'code' pe post de ID unic
      id: (json['code'] ?? '').toString(),

      // 2. Numele este 'name', deci bate perfect cu Go
      name: (json['name'] ?? 'Fără nume').toString(),

      // 3. Deoarece nu ai preț în ProductDTO, punem 0.0 deocamdată
      // (Dacă adaugi prețul mai târziu în Go, schimbi doar aici)
      price: json['price'] != null
          ? double.tryParse(json['Price'].toString()) ?? 0.0
          : 0.0,

      // 4. Folosim 'group' pe post de categorie!
      category: (json['group'] ?? 'Altele').toString(),
    );
  }
}
