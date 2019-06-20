const net = require('net');

let ready = false;

const stream = net.connect('/tmp/out')
  .on('connect', () => { ready = true; })
  .on('error', (err) => {
    console.error('IPC socket not active:', err);
    process.exit(1);
  });

module.exports = {
  get MODEL () {
    let model = {};

    try {
      model = JSON.parse(process.env.MODEL);
    } catch (e) {
      console.error('Failed to parse configuration model: ', e);
    }

    return model;
  },

  output: (message) => {
    if (!ready) {
      console.error('IPC socket not ready');
      return;
    }

    if (typeof message !== 'string') {
      try {
        message = JSON.stringify(message);
      } catch (e) {
        console.error('Failed to parse output: ', e);
        return false;
      }
    }

    length = Buffer.byteLength(message);
    // 4 bytes = 32 bits
    buffer = new Buffer(4 + length);
    // Prefix buffer length at first 32 bits
    buffer.writeUInt32BE(length, 0);
    buffer.write(message, 4);
    stream.write(buffer);
  }
};
