import 'dart:convert';
import 'package:http/http.dart' as http;
import 'auth_service.dart';
import '../core/constants.dart';
import '../core/logger.dart';
import '../models/client.dart';
import '../models/contract.dart';

class ApiService {
  // 1. CĂUTARE CLIENȚI
  Future<List<Client>> searchClients(String query) async {
    myLog("🔍 Căutăm clienți: $query");

    if (query.length < 3) return [];

    try {
      // 2. Apelăm ruta de căutare clienți
      final url = Uri.parse('${AppConfig.clientsSearchEndpoint}?q=$query');
      final response = await http.get(url, headers: AuthService.getHeaders());
      
      // 3. Logăm JSON-ul primit de la server pentru debugging
      myLog("🚀 JSON RAW DE LA SERVER: ${response.body}");

      // 4. Dacă răspunsul e OK, transformăm JSON-ul în List<Client>     
      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Client.fromJson(json)).toList();
      }
    } catch (e) {
      // 5. Dacă apare o eroare, o logăm cu myLog (nu uităm să includem și eroarea!)  
      myLog("❌ (Clients) Eroare API Căutare : $e");
    }
    return [];
  }

  // 2. FETCH CONTRACTE (Piesa care lipsea!)
Future<List<Contract>> fetchContracts(String clientId) async {
    try {
      // 1. Apelăm ruta stabilită în Gin
      final url = Uri.parse('${AppConfig.contractsEndpointByClient}/$clientId');
      final response = await http.get(url, headers: AuthService.getHeaders());

      // 2. Folosim myLog (așa cum trebuie!) pentru a vedea JSON-ul
      myLog("🚀 JSON CONTRACTE: ${response.body}");

      // 3. Dacă răspunsul e OK, transformăm JSON-ul în List<Contract>
      if (response.statusCode == 200) {
        List<dynamic> data = json.decode(response.body);
        return data.map((json) => Contract.fromJson(json)).toList();
      }
    } catch (e) {
      // 4. Dacă apare o eroare, o logăm cu myLog (nu uităm să includem și eroarea!)
      myLog("❌ Eroare la contracte", error: e);
    }
    return [];
  }
}
