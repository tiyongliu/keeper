import _compact from 'lodash-es/compact';

// original C# variant
// public bool Match(string value)
// {
//     if (String.IsNullOrEmpty(Filter)) return false;
//     if (String.IsNullOrEmpty(value)) return true;

//     var camelVariants = new HashSet<string>();
//     camelVariants.Add(new String(value.Where(Char.IsUpper).ToArray()));
//     if (value.All(x => Char.IsUpper(x) || x == '_'))
//     {
//         var sb = new StringBuilder();
//         for (int i = 0; i < value.Length; i++)
//         {
//             if (Char.IsUpper(value[i]) && (i == 0 || value[i - 1] == '_')) sb.Append(value[i]);
//         }
//         camelVariants.Add(sb.ToString());
//     }
//     else
//     {
//         string s = value, s0;
//         do
//         {
//             s0 = s;
//             s = Regex.Replace(s, "([A-Z])([A-Z])([A-Z])", "$1$3");
//         } while (s0 != s);
//         camelVariants.Add(new String(s.Where(Char.IsUpper).ToArray()));
//     }

//     bool camelMatch = camelVariants.Any(x => DoMatch(Filter, x));
//     if (Filter.All(Char.IsUpper)) return camelMatch;
//     return DoMatch(Filter, value) || camelMatch;
// }

function fuzzysearch(needle, haystack) {
  var hlen = haystack.length;
  var nlen = needle.length;
  if (nlen > hlen) {
    return false;
  }
  if (nlen === hlen) {
    return needle === haystack;
  }
  outer: for (var i = 0, j = 0; i < nlen; i++) {
    var nch = needle.charCodeAt(i);
    while (j < hlen) {
      if (haystack.charCodeAt(j++) === nch) {
        continue outer;
      }
    }
    return false;
  }
  return true;
}

export function filterName(filter: string, ...names: string[]) {
  if (!filter) return true;

  // const camelVariants = [name.replace(/[^A-Z]/g, '')]
  const tokens = filter.split(' ').map(x => x.trim());

  const namesCompacted = _compact(names);
  for (const token of tokens) {
    const tokenUpper = token.toUpperCase();
    const found = namesCompacted.find(name => fuzzysearch(tokenUpper, name.toUpperCase()));
    if (!found) return false;
  }

  return true;
}
