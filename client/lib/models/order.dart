import 'package:gorders_mobile/models/product.dart';

class OrderItem {
  final Product product;
  int quantity;

  OrderItem({required this.product, this.quantity = 0});

  // Metodă care pregătește datele exact cum le vrea Go-ul în JSON
  Map<String, dynamic> toJson() => {
    'product_id': product.id,
    'quantity': quantity,
  };
}
