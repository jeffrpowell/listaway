require("../index")

document.addEventListener('DOMContentLoaded', function () {
    const form = document.querySelector('.reset-form');
    if (!form) return;
    form.addEventListener('submit', async function (e) {
        e.preventDefault();
        const password = form.querySelector('input[name="password"]').value;
        const confirmPassword = form.querySelector('input[name="confirmPassword"]').value;
        const errorSpan = form.querySelector('.reset-form-error-span');
        errorSpan.classList.add('hidden');
        errorSpan.textContent = '';
        if (password !== confirmPassword) {
            errorSpan.textContent = 'Passwords do not match.';
            errorSpan.classList.remove('hidden');
            return;
        }
        // Get token from URL
        const token = window.location.pathname.split('/').pop();
        try {
            const res = await fetch(`/reset/${token}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                body: new URLSearchParams({ password }),
            });
            if (res.ok) {
                form.innerHTML = '<p class="text-green-600">Your password has been reset. You may now <a href="/auth" class="text-blue-500 underline">log in</a>.</p>';
            } else {
                const msg = await res.text();
                errorSpan.textContent = msg || 'Failed to reset password.';
                errorSpan.classList.remove('hidden');
            }
        } catch (err) {
            errorSpan.textContent = 'Network error.';
            errorSpan.classList.remove('hidden');
        }
    });
});
