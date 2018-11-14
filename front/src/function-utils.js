export function debounce(delay, f) {
  var id = null;
  return function () {
    var self = this;
    var args = arguments;
    clearTimeout(id);
    id = setTimeout(function () {
      f.apply(self, args);
    }, delay);
  }
}
