const btn = document.getElementById("btn");
const htmlForm = document.getElementById("url-form");
const result = document.getElementById("result");
const ok = "#58d68d";
const notok = "#f44336"

function shorten(event) {
  event.preventDefault();
  formData = new FormData(htmlForm);
  fetch('/', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(Object.fromEntries(formData)),
  }).then(res => {
    if (res.status == 201) {
      return res.json().then(data => {
        clear();
        appendText("shortened", ok);
        appendText(data.shorturl);
        if (copyToClip(data.shorturl))
          appendText("copied to clipboard!", ok);
      });
    } else {
      clear();
      appendText("server error :(", notok);
    }
  }).catch(err => {console.log(err)});
}

// Insert given string into the result
function appendText(str, color = "#171717") {
  let span = document.createElement('span');
  span.textContent = str;
  result.appendChild(span);
  span.style.color = color;
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
