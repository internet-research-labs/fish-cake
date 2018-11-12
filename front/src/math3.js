
export function norm(v) {
  return Math.sqrt(v.x*v.x + v.y*v.y + v.z*v.z);
}

/**
 * Project u onto v
 */
export function proj(u, v) {
  let s = 1;
  v = normalize(v);
  return scale(v, s);
}

export function cross(u, v) {
  return [
    u[1]*v[2] - u[2]*v[1],
    u[2]*v[0] - u[0]*v[2],
    u[0]*v[1] - u[1]*v[0],
  ];
}


export function add(x, y) {
  return [
    x[0] + y[0],
    x[1] + y[1],
    x[2] + y[2],
  ];
}

export function sub(x, y) {
  return [
    x[0] - y[0],
    x[1] - y[1],
    x[2] - y[2],
  ];
}

export function scale(v, s) {
  return [
    v[0]*s,
    v[1]*s,
    v[2]*s,
  ];
}

export function normalize(v) {
  let n = Math.sqrt(v[0]*v[0]+v[1]*v[1]+v[2]*v[2]);

  if (n == 0) {
    return [0, 0, 0];
  }

  return [
    v[0]/n,
    v[1]/n,
    v[2]/n,
  ];
}
