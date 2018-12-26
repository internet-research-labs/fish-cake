export function randomColor() {
  let r = Math.random()*255;
  let g = Math.random()*255;
  let b = Math.random()*255;
  return r*16*16*16*16 + g*16*16 + b;
}

export function makeSwift() {
  return new THREE.Mesh(
    new THREE.BoxGeometry(0.1, 0.1, 0.1),
    new THREE.MeshBasicMaterial({color: randomColor()}),
  );
}
