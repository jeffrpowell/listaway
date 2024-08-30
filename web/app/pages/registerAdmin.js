require("../index")

async function sendData(form) {
    // Associate the FormData object with the form element
    const formData = new FormData(form);

    try {
        const response = await fetch("/admin/register", {
            method: "PUT",
            headers: {
                "Accept": "text/plain",
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(formData).toString()
        });
        if (response.status >= 400) {
            showError(response);
        } else if (response.status === 200) {
            window.location.href = response.headers.get("Location");
        }
    } catch (e) {
        console.error(e);
        showError(500);
    }
}

// Take over form submission
const forms = document.querySelectorAll(".register-admin-form");
forms.forEach(form => form.addEventListener("submit", (event) => {
    event.preventDefault();
    sendData(form);
}));

async function showError(response) {
    text = await response.text();
    const errorSpans = document.querySelectorAll(".error-span");
    errorSpans.forEach(errorSpan => {
        if (response.status < 500) {
            errorSpan.innerText = text;
        }
        else {
            errorSpan.innerText = "Unexpected error occurred. Please try again later.";
        }
        errorSpan.classList.remove("hidden");
    });
}