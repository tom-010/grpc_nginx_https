import 'package:client/channel_builders/build_channel_native.dart' if(dart.library.js) 'package:client/channel_builders/build_channel_web.dart';
import 'package:grpc/grpc_connection_interface.dart';

ClientChannelBase bulidChannel() {
  return bulidChannelImpl();
}