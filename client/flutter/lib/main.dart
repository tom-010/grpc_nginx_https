import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:grpc/grpc.dart';
import 'package:grpc/grpc_web.dart';

import 'proto/greeter.pbgrpc.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({Key? key}) : super(key: key);

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(title: 'Flutter gRPC client'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({Key? key, required this.title}) : super(key: key);
  final String title;
  @override
  State<MyHomePage> createState() => _MyHomePageState();
}


class _MyHomePageState extends State<MyHomePage> {
  List<String> messages = [];
  bool loading = false;

  void log(String msg) {
    setState(() {
      messages.add(msg);
    });
  }

  void _send() async {
    setState(() {
      loading = true;
    });

    
    final channel = kIsWeb
        ? GrpcWebClientChannel.xhr(Uri.parse('https://web-api.example.com:443')) 
        : ClientChannel('api.example.com',
            port: 443,
            options: ChannelOptions(
              credentials: const ChannelCredentials.secure(),
              codecRegistry:
                  CodecRegistry(codecs: const [GzipCodec(), IdentityCodec()]),
            ));


    String token = base64Encode(utf8.encode('tom:password'));
    final client = GreeterClient(channel, options: CallOptions(metadata: {'authorization': 'Basic $token'}));

    try {
      final response = await client.sayHello(HelloRequest(name: 'Tom'));
      log(response.message);
    } catch (e) {
      log('Cought error: $e');
      setState(() {
        loading = false;
      });
    }

    ResponseStream<HelloReply> stream =
        client.sayRepeatHello(RepeatHelloRequest(name: 'Tom', count: 7));

    await for (HelloReply r in stream) {
      log(r.message);
    }

    await channel.shutdown();

    setState(() {
      loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: ListView.builder(
        itemCount: messages.length,
        itemBuilder: (BuildContext context, int index) {
          return ListTile(title: Text(messages[index]));
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: loading ? null : _send,
        tooltip: 'Send',
        backgroundColor: loading ? Colors.grey : Colors.blueAccent,
        child: const Icon(Icons.send),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
