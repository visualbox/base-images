from os import environ, path
from json import loads, dumps
from socket import socket, AF_UNIX, SOCK_STREAM
from ctypes import c_int32
from struct import pack

if path.exists("/tmp/out"):
  stream = socket(AF_UNIX, SOCK_STREAM)
  stream.connect("/tmp/out")
else:
  stream = None
  print("IPC socket not active")

MODEL = {}
try:
  MODEL = loads(environ['MODEL'])
except Exception:
  print("Failed to parse configuration model")

def output(message):
  if stream is None:
    print("IPC socket not active")
    return
  
  if not isinstance(message, str):
    try:
      message = dumps(message)
    except Exception:
      print("failed to parse output")
      return

  msgBytes = message.encode("utf-8")
  length = len(msgBytes)
  headerBuf = pack(">I", length)
  stream.send(headerBuf + msgBytes)
