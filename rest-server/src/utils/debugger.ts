import * as fs from 'fs';
import * as path from 'path';

const saveProtoAsJSON = async (filename: string, obj: any) => {
  const debugFolder = 'debug';
  const hasFolder = fs.existsSync(debugFolder);
  if (!hasFolder) {
    fs.mkdirSync(debugFolder);
  }
  try {
    fs.writeFileSync(path.join(debugFolder, filename), obj.encodeJSON());
  } catch (e) {
    console.debug(`Could not write file ${filename}, error: ${e}`);
    throw e;
  }
};

export { saveProtoAsJSON };
