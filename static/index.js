const form = document.getElementById("url-checker");
const button = document.querySelector("button[type='submit']");

form.addEventListener("submit", (event) => {
  event.preventDefault();

  const url = new FormData(form).get("url");
  if (!url) {
    alert("Please enter a URL");
    return;
  }

  if (!URL.canParse(url)) {
    alert("Invalid URL");
    return;
  }

  parsedUrl = new URL(url);
  if (!parsedUrl.protocol.startsWith("http")) {
    alert("Invalid URL");
    return;
  }
  if (!parsedUrl.hostname) {
    alert("Invalid URL");
    return;
  }

  event.currentTarget.reset();

  button.disabled = true;
  const prompt = document.createElement("strong");
  prompt.textContent = "Checking...";
  document.body.appendChild(prompt);

  fetch("/check", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ url }),
  })
    .then(async (response) => {
      if (!response.ok) {
        const errorText = await response.text().catch(() => "Unknown error");
        throw new Error(`HTTP error! status: ${response.status}, ${errorText}`);
      }
      return response.json();
    })
    .then((data) => {
      console.log(data);
      const strong = document.createElement("strong");
      strong.textContent = `Found ${data.deadLinks.length} dead links`;

      const ul = document.createElement("ul");
      data.deadLinks.forEach((link) => {
        const li = document.createElement("li");
        li.textContent = link;
        ul.appendChild(li);
      });

      document.body.appendChild(strong);
      document.body.appendChild(ul);
    })
    .catch((error) => alert(`Error fetching URL: ${error}`))
    .finally(() => {
      button.disabled = false;
      document.body.removeChild(prompt);
    });
});
