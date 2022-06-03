import 'package:grpc/grpc_connection_interface.dart';
import 'package:grpc/grpc_web.dart';

ClientChannelBase bulidChannelImpl() {
  return GrpcWebClientChannel.xhr(Uri.parse('https://web-api.example.com'));
}
