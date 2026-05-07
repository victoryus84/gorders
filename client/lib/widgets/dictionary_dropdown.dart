import 'package:flutter/material.dart';
import 'package:multi_dropdown/multi_dropdown.dart';

/// Универсальный виджет для выбора из справочника
class DictionaryDropdown<T extends Object> extends StatelessWidget {
  final String hintText; // подсказка (например: "Контрагент", "Номенклатура")
  final List<T> items; // список любых объектов
  final MultiSelectController<T> controller; // контроллер
  final String Function(T) labelBuilder; // как получить имя для отображения
  final bool singleSelect; // выбор одного или многих
  final String? Function(List<DropdownItem<T>>?)? validator; // валидатор

  const DictionaryDropdown({
    super.key,
    required this.hintText,
    required this.items,
    required this.controller,
    required this.labelBuilder,
    this.singleSelect = true,
    this.validator,
  });

  @override
  Widget build(BuildContext context) {
    return MultiDropdown<T>(
      singleSelect: singleSelect,
      items: items
          .map(
            (item) => DropdownItem<T>(label: labelBuilder(item), value: item),
          )
          .toList(),
      controller: controller,
      searchEnabled: true,
      fieldDecoration: FieldDecoration(
        hintText: hintText,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
      ),
      validator: validator,
      onSelectionChange: (selectedItems) {
        debugPrint("Выбрано: $selectedItems");
      },
    );
  }
}
