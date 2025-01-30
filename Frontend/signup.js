document.addEventListener('DOMContentLoaded', () => {
    console.log('Signup page fully loaded'); // Log when the DOM is ready

    // Event listener for the signup button
    const signupButton = document.getElementById('signupButton');
    if (signupButton) {
        signupButton.addEventListener('click', () => {
            const firstName = document.getElementById('signupFirstName').value;
            const lastName = document.getElementById('signupLastName').value;
            const email = document.getElementById('signupEmail').value;
            const password = document.getElementById('signupPassword').value;
            signupUser(firstName, lastName, email, password); // Call the signup function
        });
    } else {
        console.error('Signup button not found');
    }
});

// Function to handle user signup
const signupUser = (firstName, lastName, email, password) => {
    console.log('Attempting to sign up with email:', email); // Log the email being used for signup
    fetch('http://localhost:8000/user/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ firstName, lastName, email, password }),
    })
    .then(response => {
        console.log('Signup response status:', response.status); // Log the response status
        if (!response.ok) {
            throw new Error('Signup failed');
        }
        return response.json();
    })
    .then(data => {
        console.log('Signup successful:', data); // Log the successful signup response
        document.getElementById('signupMessage').innerText = 'Signup successful! You can now log in.';
        window.location.href = 'index.html'; // Redirect to login page after successful signup
    })
    .catch(error => {
        document.getElementById('signupMessage').innerText = error.message;
        console.error('Error signing up:', error);
    });
};