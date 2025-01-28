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

function displayResults(movies) {
    const resultsDiv = document.getElementById('results');
    resultsDiv.innerHTML = '';

    movies.forEach(movie => {
        const movieDiv = document.createElement('div');
        movieDiv.className = 'movie-item';
        movieDiv.innerText = movie.title;
        movieDiv.onclick = () => playTrailer(movie._id);
        resultsDiv.appendChild(movieDiv);
    });
}

