import 'package:grpc/grpc.dart';
import 'package:grpc/grpc_connection_interface.dart';

ClientChannelBase bulidChannelImpl() {
  return ClientChannel('api.example.com',
      port: 443,
      options: ChannelOptions(
        credentials: const ChannelCredentials.secure(),
        codecRegistry:
            CodecRegistry(codecs: const [GzipCodec(), IdentityCodec()]),
      ));
}
