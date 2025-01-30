document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM fully loaded and parsed'); // Log when the DOM is ready

    // Function to handle user login
    const loginUser  = (email, password) => {
        console.log('Attempting to log in with email:', email); // Log the email being used for login
        fetch('http://localhost:8000/user/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password }),
        })
        .then(response => {
            console.log('Login response status:', response.status); // Log the response status
            if (!response.ok) {
                throw new Error('Login failed');
            }
            return response.json();
        })
        .then(data => {
            localStorage.setItem('jwtToken', data.token);
            console.log('Login successful, token stored:', data.token); // Log the stored token
            document.getElementById('loginMessage').innerText = 'Login successful!';
            window.location.href = 'main.html';
        })
        .catch(error => {
            document.getElementById('loginMessage').innerText = error.message;
            console.error('Error logging in:', error);
        });
    };

    // Event listeners for signup and login buttons
    const signupButton = document.getElementById('signupButton');
    const loginButton = document.getElementById('loginButton');

    if (signupButton) {
        signupButton.addEventListener('click', () => {
            const email = document.getElementById('signupEmail').value;
            const password = document.getElementById('signupPassword').value;
            signupUser (email, password);
        });
    } else {
        console.error('Signup button not found');
    }

    if (loginButton) {
        loginButton.addEventListener('click', () => {
            const email = document.getElementById('loginEmail').value;
            const password = document.getElementById('loginPassword').value;
            loginUser (email, password);
        });
    } else {
        console.error('Login button not found');
    }

    // Function to handle searching for movies
    const searchInput = document.getElementById('search');
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const query = this.value;
            console.log('Search query:', query); // Log the search query

            if (query.length > 2) {
                const token = localStorage.getItem('jwtToken');
                console.log('Token being sent:', token); // Log the token

                fetch(`http://localhost:8000/search?q=${query}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': token, // Send the token directly without "Bearer"
                    },
                })
                .then(response => {
                    console.log('Search response status:', response.status); // Log the response status
 if (!response.ok) {
                        throw new Error(`HTTP error! Status: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    console.log('Search results:', data); // Log the search results
                    displayResults(data);
                })
                .catch(error => console.error('Error fetching movies:', error));
            } else {
                document.getElementById('results').innerHTML = ''; // Clear results if query is too short
            }
        });
    } else {
        console.error('Search input not found');
    }

    // Function to display search results
    function displayResults(movies) {
        const resultsDiv = document.getElementById('results');
        resultsDiv.innerHTML = ''; // Clear previous results

        if (movies.length === 0) {
            resultsDiv.innerHTML = '<p>No results found.</p>'; // Show message if no results
            return;
        }

        movies.forEach(movie => {
            const movieDiv = document.createElement('div');
            movieDiv.className = 'movie-item';
            movieDiv.innerText = movie.title; // Display the movie title
            movieDiv.onclick = () => showMovieDetails(movie); // Show details on click
            resultsDiv.appendChild(movieDiv); // Append the movie title to results
        });
    }

  // Function to show movie details
function showMovieDetails(movie) {
    console.log('Showing details for movie:', movie.title); // Log the movie title being shown
    document.getElementById('movieTitle').innerText = movie.title; // Set the movie title
    document.getElementById('moviePlot').innerText = movie.plot; // Set the movie plot
    document.getElementById('movieYear').innerText = movie.year; // Set the movie year
    document.getElementById('movieGenres').innerText = movie.genres.join(', '); // Set the movie genres
    document.getElementById('movieRating').innerText = movie.imdb.rating; // Set the movie rating
    document.getElementById('movieDetails').style.display = 'block'; // Show movie details

    // Set up the play trailer button
    const playTrailerButton = document.getElementById('playTrailerButton');
    if (playTrailerButton) {
        playTrailerButton.onclick = () => {
            playTrailer(movie._id); // Call the function to play the trailer
        };
    } else {
        console.error('Play trailer button not found');
    }
}

    // Function to handle playing the trailer
    const playTrailer = (id) => {
        const token = localStorage.getItem('jwtToken'); // Retrieve the token
        console.log('Fetching trailer for movie ID:', id); // Log the movie ID being fetched

        fetch(`http://localhost:8000/movie/${id}/playback`, {
            method: "GET",
            headers: {
                'Authorization': token, // Include the token in the Authorization header
            },
        })
        .then((response) => {
            console.log('Trailer fetch response status:', response.status); // Log the response status
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then((data) => {
            console.log('Trailer data received:', data); // Log the trailer data received
            if (!data.trailer_url) {
                throw new Error("Trailer URL not found in response");
            }

            const trailerIframe = document.getElementById("trailerIframe");
            const trailerModal = document.getElementById("trailerModal");

            if (trailerIframe) {
                trailerIframe.src = data.trailer_url; // Set the iframe source to the trailer URL
            } else {
                console.error('Trailer iframe not found');
            }

            if (trailerModal) {
                trailerModal.style.display = "flex"; // Show the modal
            } else {
                console.error('Trailer modal not found');
            }
        })
        .catch((error) => {
            console.error("Error fetching trailer URL:", error);
            alert("Unable to load the trailer. Please try again later.");
        });
    };

    // Function to close the trailer modal
    const closeTrailer = () => {
        const trailerModal = document.getElementById("trailerModal");
        const trailerIframe = document.getElementById("trailerIframe");

        if (trailerIframe) {
            trailerIframe.src = ""; // Clear the iframe source to stop the video
        } else {
            console.error('Trailer iframe not found');
        }

        if (trailerModal) {
            trailerModal.style.display = "none"; // Hide the modal
        } else {
            console.error('Trailer modal not found');
        }
        console.log('Trailer modal closed'); // Log when the modal is closed
    };

    const closeTrailerButton = document.getElementById("closeTrailerButton");
if (closeTrailerButton) {
    closeTrailerButton.onclick = closeTrailer;
} else {
    console.error('Close trailer button not found');
}
});