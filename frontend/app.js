input = document.querySelector("#search-input");
searchBtn = document.querySelector("#search-button");
suggestedDiv = document.querySelector("#suggested-results");

function fireSearch(value, confirmed) {
  return fetch("http://127.0.0.1:8090/search", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ prefix: value, confirmed: confirmed })
  })
    .then(response => response.json())
    .then(result => updateResult(result))
    .catch(err => console.log(err.message));
}

function updateResult(words) {
  while (suggestedDiv.firstChild) {
    suggestedDiv.removeChild(suggestedDiv.firstChild);
  }
  if (words) {
    words.forEach(word => {
      const p = document.createElement("p");
      p.innerText = word;
      suggestedDiv.appendChild(p);
    });
  }
}

input.addEventListener("input", e => {
  const { value } = e.target;
  if (value.length > 0) {
    fireSearch(value, false);
  } else {
    updateResult(null);
  }
});

input.addEventListener("keyup", e => {
  if (e.keyCode === 13) {
    const { value } = e.target;
    fireSearch(value, true);
  }
});

searchBtn.addEventListener("click", e => {
  const { value } = e.target;
  if (value.length > 0) {
    fireSearch(value, true);
  }
});
var myHeaders = new Headers();
