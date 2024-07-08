require("../index")

async function sendData(form) {
    // Associate the FormData object with the form element
    const formData = new FormData(form);

    try {
        const response = await fetch("/auth", {
            method: "POST",
            // Set the FormData instance as the request body
            headers: {
                "Accept": "text/plain",
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(formData).toString()
        });
        if (response.status >= 400) {
            showError(response.status);
        } else if (response.status === 200) {
            window.location.href = response.headers.get("Location");
        }
    } catch (e) {
        console.error(e);
        showError(500);
    }
}

// Take over form submission
const forms = document.querySelectorAll(".login-form");
forms.forEach(form => form.addEventListener("submit", (event) => {
    event.preventDefault();
    sendData(form);
}));

function showError(statusCode) {
    const errorSpans = document.querySelectorAll(".error-span");
    errorSpans.forEach(errorSpan => {
        if (statusCode === 401) {
            errorSpan.innerText = "Email or password is not recognized.";
        }
        else {
            errorSpan.innerText = "Unexpected error occurred. Please try again later.";
        }
        errorSpan.classList.remove("hidden");
    });
}