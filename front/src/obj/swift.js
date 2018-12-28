export function randomGreyColor() {
  let l = Math.floor(Math.random()*256);
  let c = l*0x00FFFF + l*0x0000FF + l;
  return l*256*256 + l*256 + l;
}

export function randomColor() {
  let r = Math.random()*255;
  let g = Math.random()*255;
  let b = Math.random()*255;
  return r*16*16*16*16 + g*16*16 + b;
}

export function makeSwift() {
  let size = 0.2;
  let m = new THREE.Mesh(
    new THREE.BoxGeometry(size/4.0, size/4.0, 2*size),
    new THREE.MeshBasicMaterial({color: randomGreyColor()}),
  );
  m.rotation.z = -Math.PI/2.0;
  return m;
}
