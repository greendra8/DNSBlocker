<script>
  import Nav from '$lib/components/Nav.svelte';
</script>

<Nav />

<main>
  <h1>Welcome to DNS Manager</h1>
  <p>Manage your internet habits and boost productivity with our DNS-based blocking system.</p>
  <div class="cta-buttons">
    <a href="/habits" class="button">Track Habits</a>
    <a href="/block-sites" class="button">Manage Blocked Sites</a>
  </div>
</main>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&family=Roboto:wght@300;400&display=swap');

  :global(body) {
    font-family: 'Roboto', sans-serif;
    color: #FFFFFF;
    background-color: #000033;
    margin: 0;
    padding: 0;
  }

  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem;
    text-align: center;
  }

  h1 {
    font-family: 'Playfair Display', serif;
    font-size: 2.5rem;
    color: #FFFFFF;
    margin-bottom: 1rem;
  }

  p {
    font-size: 1.2rem;
    margin-bottom: 2rem;
    line-height: 1.4;
  }

  .cta-buttons {
    display: flex;
    justify-content: center;
    gap: 20px;
  }

  .button {
    display: inline-block;
    padding: 0.8rem 1.5rem;
    background-color: #FFFFFF;
    color: #000033;
    text-decoration: none;
    font-size: 1rem;
    font-weight: bold;
    transition: background-color 0.3s ease;
  }

  .button:hover {
    background-color: #CCCCCC;
  }
</style>