require("../index")

async function sendData(form) {
    const formData = new FormData(form);
    const mode = form.dataset.mode;              // "login" or "reset"
    const endpoint = mode === "reset" ? "/reset" : "/auth";

    // if we're resetting, we don't want the password field at all:
    if (mode === "reset") {
        formData.delete("password");
    }

    try {
        const response = await fetch(endpoint, {
            method: "POST",
            headers: {
                "Accept": "text/plain",
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: new URLSearchParams(formData).toString()
        });
        if (response.status >= 400) {
            showError(response.status);
        } else if (response.status === 200) {
            if (mode === "login") {
                // normal sign‑in: redirect to Location header
                window.location.href = response.headers.get("Location");
            } else {
                // reset requested – show a success message
                alert("If that email is in our system you’ll receive reset instructions shortly.");
                // Optional: go back to login form automatically
                showLoginForm();
            }
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

// cache a few elements
const pwdGroup = document.getElementById("password-group");
const pwdInput = document.getElementById("password");
const submitBtn = document.getElementById("submit-btn");
const forgotLink = document.getElementById("forgot-link");
const errorSpan = document.querySelector(".error-span");

// Show the reset‑password “mode”
function showResetForm() {
    form.dataset.mode = "reset";
    pwdGroup.classList.add("hidden");      // hide the password field
    pwdInput.disabled = true;              // ensure it's not sent
    submitBtn.innerText = "Request Reset"; // swap button text
    forgotLink.innerText = "Back to sign in form";
}

// Show the normal login “mode”
function showLoginForm() {
    form.dataset.mode = "login";
    pwdGroup.classList.remove("hidden");
    pwdInput.disabled = false;
    submitBtn.innerText = "Sign In";
    forgotLink.innerText = "Forgot password?";
    errorSpan.classList.add("hidden");     // clear any lingering error
}

// Toggle when someone clicks the link
forgotLink.addEventListener("click", e => {
    e.preventDefault();
    if (form.dataset.mode === "login") {
        showResetForm();
    } else {
        showLoginForm();
    }
});

// existing error display (unchanged)
function showError(statusCode) {
    if (statusCode === 401) {
        errorSpan.innerText = "Email or password is not recognized.";
    } else {
        errorSpan.innerText = "Unexpected error occurred. Please try again later.";
    }
    errorSpan.classList.remove("hidden");
}
