require("../index")

const pwdGroup = document.getElementById("password-group");
const pwdInput = document.getElementById("password");
const submitBtn = document.getElementById("submit-btn");
const forgotLink = document.getElementById("forgot-link");
const errorSpan = document.getElementById("error-span");
const oidcSection = document.getElementById("oidc-section");
const oidcLoginBtn = document.getElementById("oidc-login-btn");
const oidcProviderText = document.getElementById("oidc-provider-text");

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
                showError(200);
                showLoginForm();
            }
        }
    } catch (e) {
        console.error(e);
        showError(500);
    }
}

// Take over form submission
const form = document.getElementById("login-form");
form.addEventListener("submit", e => {
    e.preventDefault();
    sendData(form);
});

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

// existing error display
function showError(statusCode) {
    if (statusCode === 401) {
        errorSpan.innerText = "Email or password is not recognized.";
    } else if (statusCode === 200){
        errorSpan.innerText = "If that email is in our system, it will receive reset instructions shortly."
    } else {
        errorSpan.innerText = "Unexpected error occurred. Please try again later.";
    }
    errorSpan.classList.remove("hidden");
}

// OIDC functionality
async function checkOIDCStatus() {
    try {
        const response = await fetch("/api/oidc/status");
        if (response.ok) {
            const data = await response.json();
            if (data.enabled) {
                oidcSection.classList.remove("hidden");
                
                // Update button text based on provider
                if (data.provider) {
                    const providerName = data.provider.charAt(0).toUpperCase() + data.provider.slice(1);
                    oidcProviderText.innerText = `Continue with ${providerName}`;
                }
            }
        }
    } catch (e) {
        console.log("OIDC status check failed:", e);
        // OIDC not available, keep section hidden
    }
}

// Handle OIDC login
function handleOIDCLogin() {
    // Redirect to OIDC login endpoint
    window.location.href = "/auth/oidc/login";
}

// Add OIDC event listener
if (oidcLoginBtn) {
    oidcLoginBtn.addEventListener("click", handleOIDCLogin);
}

// Check URL parameters for OIDC errors
function checkOIDCErrors() {
    const urlParams = new URLSearchParams(window.location.search);
    const error = urlParams.get('error');
    
    if (error === 'oidc_failed') {
        showError(500);
        errorSpan.innerText = "Authentication failed. Please try again.";
    }
}

// Initialize OIDC functionality
document.addEventListener("DOMContentLoaded", () => {
    checkOIDCStatus();
    checkOIDCErrors();
});
