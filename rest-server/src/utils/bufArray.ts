const ab2str = (buffer: Uint8Array): String => {
  return String.fromCharCode.apply(null, new Uint8Array(buffer));
};

const str2ab = (str: String): ArrayBuffer => {
  const buf = new ArrayBuffer(str.length * 2);
  let bufView = new Uint8Array(buf);
  [...str].forEach((element, index) => {
    bufView[index] = element.charCodeAt(index);
  });
  return buf;
};

export { ab2str, str2ab };
