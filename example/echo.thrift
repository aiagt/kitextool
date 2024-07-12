namespace go echo

service EchoService {
    string Echo(1: string message)
}