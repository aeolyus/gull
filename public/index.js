const ok = "#58d68d";
const notok = "#f44336";

const htmlForm = document.getElementById("url-form");
const url = document.getElementById("url");
const alias = document.getElementById("alias");
const btn = document.getElementById("btn");
const result = document.getElementById("result");

// RFC 3986 Section 2.3 URI Unreserved Characters
const URIUnreservedChars = /^([A-Za-z0-9_.~-])+$/;
const aliasRegex = RegExp(URIUnreservedChars);
// https://gist.github.com/dperini/729294
const URLRegex = RegExp(
  "^" +
    // protocol identifier (optional)
    // short syntax // still required
    "(?:(?:(?:https?|ftp):)?\\/\\/)" +
    // user:pass BasicAuth (optional)
    "(?:\\S+(?::\\S*)?@)?" +
    "(?:" +
      // IP address dotted notation octets
      // excludes loopback network 0.0.0.0
      // excludes reserved space >= 224.0.0.0
      // excludes network & broadcast addresses
      // (first & last IP address of each class)
      "(?:[1-9]\\d?|1\\d\\d|2[01]\\d|22[0-3])" +
      "(?:\\.(?:1?\\d{1,2}|2[0-4]\\d|25[0-5])){2}" +
      "(?:\\.(?:[1-9]\\d?|1\\d\\d|2[0-4]\\d|25[0-4]))" +
    "|" +
      // host & domain names, may end with dot
      // can be replaced by a shortest alternative
      // (?![-_])(?:[-\\w\\u00a1-\\uffff]{0,63}[^-_]\\.)+
      "(?:" +
        "(?:" +
          "[a-z0-9\\u00a1-\\uffff]" +
          "[a-z0-9\\u00a1-\\uffff_-]{0,62}" +
        ")?" +
        "[a-z0-9\\u00a1-\\uffff]\\." +
      ")+" +
      // TLD identifier name, may end with dot
      "(?:[a-z\\u00a1-\\uffff]{2,}\\.?)" +
    ")" +
    // port number (optional)
    "(?::\\d{2,5})?" +
    // resource path (optional)
    "(?:[/?#]\\S*)?" +
  "$"
);

function shorten(event) {
  clear();
  event.preventDefault();
  if (!isValidURL(url.value)) {
    appendText("invalid url", notok);
    return
  }
  if (alias.value != "" && !isValidAlias(alias.value)) {
    appendText("invalid alias, use only alphanumerics, underscores, periods, tildes, and dashes", notok);
    return
  }
  sendReq()
}

function sendReq() {
  formData = new FormData(htmlForm);
  fetch('/', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(Object.fromEntries(formData)),
  }).then(res => {
    if (res.status == 201) {
      return res.json().then(data => {
        appendText("shortened", ok);
        appendText(data.shorturl);
        if (copyToClip(data.shorturl))
          appendText("copied to clipboard!", ok);
      });
    } else {
      appendText("server error :(", notok);
    }
  }).catch(err => {console.log(err)});
}

function isValidURL(str) {
  return URLRegex.test(str)
}

// Verify alias contains only URI unreserved chars
function isValidAlias(str) {
  return aliasRegex.test(str)
}

// Insert given string into the result
function appendText(str, color = "") {
  let span = document.createElement('span');
  span.textContent = str;
  result.appendChild(span);
  if (color) {
    span.style.color = color;
  }
}

// Clear out the results
function clear() {
  result.innerHTML = "";
}

// Copy to clipboard
function copyToClip(str) {
  return navigator.clipboard.writeText(str).then(() => true, () => false);
}

btn.addEventListener("click", shorten)
