import socket

def encode_varint(n):
    while True:
        m = n & 0x7f
        n >>= 7
        if n != 0:
            yield m | 0x80
        else:
            yield m
            break

def decode_varint(iterable):
    shift = 0
    n = 0
    for byte in iterable:
        n |= (byte & 0x7f) << shift
        shift += 7
        if byte & 0x80 == 0:
            break
    else: # didn't see a terminating byte
        raise EOFError

    return n

def socket_byte_iterator(reader):
    while True:
        byte = reader.recv(1)
        if not byte:
            break
        yield ord(byte)

def rpc(address, request, response):
    with socket.create_connection(address) as conn:
        # send request
        send_bytes = request.SerializeToString()
        length = len(send_bytes)
        conn.send(bytes(encode_varint(length)))
        conn.send(send_bytes)
        # recv response
        length = decode_varint(socket_byte_iterator(conn))
        recv_bytes = conn.recv(length)
        response.ParseFromString(recv_bytes)
        return response

__ALL__ = ['rpc']

# vim: ts=4:sw=4:et
