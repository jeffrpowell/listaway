require('../navbar');

document.addEventListener('DOMContentLoaded', (event) => {
    const collectionNameInputs = document.getElementsByName('name');
    const submitButtons = document.querySelectorAll('button');
    const errorSpans = document.querySelectorAll('.error-span');
    var formReadyToSubmit = false;

    collectionNameInputs.forEach(collectionNameInput => 
        collectionNameInput.addEventListener('input', (event) => {
            submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
            errorSpans.forEach(errorSpan => {
                errorSpan.classList.add('hidden');
                errorSpan.textContent = '';
            });
            if (collectionNameInput.value.trim() !== '') {
                debouncedCheckCollectionName(submitButtons, errorSpans, collectionNameInput.value.trim());
            }
            else {
                formReadyToSubmit = false;
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
            }
        })
    );

    var debouncedCheckCollectionName = debounce(checkCollectionName, 500);
    async function checkCollectionName(submitButtons, errorSpans, name) {
        try {
            const response = await fetch(`/collections/namecheck?name=${encodeURIComponent(name)}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.status === 200) {
                submitButtons.forEach(submitButton => submitButton.classList.remove("opacity-50", "cursor-not-allowed"));
                formReadyToSubmit = true;
            } else if (response.status === 400 || response.status === 409) {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'The collection name is already in use.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            } else {
                submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
                errorSpans.forEach(errorSpan => {
                    errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                    errorSpan.classList.remove('hidden');
                });
                formReadyToSubmit = false;
            }
        } catch (error) {
            submitButtons.forEach(submitButton => submitButton.classList.add("opacity-50", "cursor-not-allowed"));
            errorSpans.forEach(errorSpan => {
                errorSpan.textContent = 'There was a problem while checking if the name was taken. Please try again.';
                errorSpan.classList.remove('hidden');
            });
            formReadyToSubmit = false;
        }
    }
    
    async function sendData(form) {
        if (!formReadyToSubmit) {
            return;
        }
        // Associate the FormData object with the form element
        const formData = new FormData(form);

        try {
            const response = await fetch("/collections", {
                method: "POST",
                headers: {
                    "Accept": "text/plain",
                    "Content-Type": "application/x-www-form-urlencoded"
                },
                body: new URLSearchParams(formData).toString()
            });
            if (response.status >= 400) {
                showError(response.status);
            } else if (response.status === 200 || response.status === 201) {
                const location = response.headers.get("Location") || `/collections/${await response.text()}`;
                window.location.href = location;
            }
        } catch (e) {
            console.error(e);
            showError(500);
        }
    }

    // Take over form submission
    const forms = document.querySelectorAll(".collection-form");
    forms.forEach(form => form.addEventListener("submit", (event) => {
        event.preventDefault();
        sendData(form);
    }));

    function showError(statusCode) {
        const errorSpans = document.querySelectorAll(".error-span");
        errorSpans.forEach(errorSpan => {
            if (statusCode === 400 || statusCode === 409) {
                errorSpan.innerText = "Collection name is already in use.";
            }
            else {
                errorSpan.innerText = "Unexpected error occurred. Please try again later.";
            }
            errorSpan.classList.remove("hidden");
        });
    }

    function debounce(func, delay) {
        let timeoutId;
        return function(...args) {
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
            timeoutId = setTimeout(() => {
                func.apply(this, args);
            }, delay);
        };
    }
});
