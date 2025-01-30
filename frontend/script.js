// Search bar event listener
document.getElementById('search').addEventListener('input', function() {
    const query = this.value;
    if (query.length > 2) {
        fetch(`http://localhost:8000/search?q=${query}`)
            .then(response => response.json())
            .then(data => {
                displayResults(data);
            })
            .catch(error => console.error('Error fetching movies:', error));
    } else {
        document.getElementById('results').innerHTML = '';
    }
});

// Function to display search results
function displayResults(movies) {
    const resultsDiv = document.getElementById('results');
    resultsDiv.innerHTML = '';

    movies.forEach(movie => {
        const movieDiv = document.createElement('div');
        movieDiv.className = 'movie-item';
        movieDiv.innerText = movie.title;
        movieDiv.onclick = () => playTrailer(movie._id); // Attach click event to play trailer
        resultsDiv.appendChild(movieDiv);
    });
}

// Function to handle playing the trailer
const playTrailer = (id) => {
    // Make the fetch call with the correct dynamic URL
    fetch(`http://localhost:8000/movie/${id}/playback`, {
        method: "GET",
    })
        .then((response) => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then((data) => {
            if (!data.trailer_url) {
                throw new Error("Trailer URL not found in response");
            }

            console.log("Trailer URL:", data.trailer_url);

            // Open the YouTube trailer in an iframe or a modal
            const trailerIframe = document.getElementById("trailerIframe");
            const trailerModal = document.getElementById("trailerModal");

            // Set the iframe source to the trailer URL
            trailerIframe.src = data.trailer_url;

            // Display the modal
            trailerModal.style.display = "block";
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

    // Stop the video when closing the modal
    trailerIframe.src = "";
    trailerModal.style.display = "none";
};

// Attach event listeners to all movies dynamically
document.addEventListener("DOMContentLoaded", () => {
    const movieElements = document.querySelectorAll(".movie");
    movieElements.forEach((movieElement) => {
        const movieId = movieElement.getAttribute("data-id");
        movieElement.onclick = () => playTrailer(movieId);
    });

    // Attach close button functionality
    const closeButton = document.getElementById("closeTrailerButton");
    closeButton.onclick = closeTrailer;
});
