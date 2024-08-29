<script>
    import { onMount } from 'svelte';
    import Nav from '$lib/components/Nav.svelte';
    import axios from 'axios';
  
    /** @type {string} User ID input */
    let userId = '';
    /** @type {string} Domain input */
    let domain = '';
    /** @type {string} Message to display to the user */
    let message = '';
    /** @type {Array<{id: number, domain: string}>} List of blocked domains */
    let rules = [];
  
    /** @type {string} Base URL for API requests */
    const API_BASE_URL = 'https://localhost:443'; // Update this to your CoreDNS server address
  
    // Load userId from localStorage on mount
    onMount(() => {
      const storedUserId = localStorage.getItem('userId');
      if (storedUserId) {
        userId = storedUserId;
        fetchRules();
      }
    });
  
    // Save userId to localStorage when it changes
    $: {
      if (userId) {
        localStorage.setItem('userId', userId);
      }
    }
  
    /**
     * Adds a new rule to block a domain
     * @returns {Promise<void>}
     */
    async function addRule() {
      if (!userId || !domain) {
        message = 'Please enter both User ID and Domain';
        return;
      }
      try {
        await axios.post(`${API_BASE_URL}/add_rule`, {
          user_id: userId,
          domain: domain
        });
        message = `Rule added successfully: ${domain}`;
        domain = '';
        await fetchRules();
      } catch (error) {
        message = 'Error adding rule: ' + error.response?.data || error.message;
      }
    }
  
    /**
     * Removes a rule
     * @param {number} ruleId - The ID of the rule to remove
     * @param {string} domainName - The domain name of the rule being removed
     * @returns {Promise<void>}
     */
    async function removeRule(ruleId, domainName) {
      try {
        await axios.post(`${API_BASE_URL}/remove_rule`, {
          user_id: userId,
          rule_id: ruleId
        });
        message = `Rule removed successfully: ${domainName}`;
        await fetchRules();
      } catch (error) {
        message = 'Error removing rule: ' + error.response?.data || error.message;
      }
    }
  
    /**
     * Fetches rules for the current user
     * @returns {Promise<void>}
     */
    async function fetchRules() {
      if (!userId) {
        message = 'Please enter a user ID';
        rules = [];
        return;
      }
      try {
        const response = await axios.get(`${API_BASE_URL}/rules/${userId}`);
        rules = response.data;
        if (rules.length === 0) {
          message = 'No rules found for this user.';
        } else {
          message = `Found ${rules.length} rule(s) for user ${userId}`;
        }
      } catch (error) {
        message = 'Error fetching rules: ' + error.response?.data || error.message;
        rules = [];
      }
    }
  </script>
  
  <Nav />

  <main>
    <h1>Manage Blocked Sites</h1>
    <form on:submit|preventDefault={addRule}>
      <label>
        User ID:
        <input bind:value={userId} on:change={fetchRules} required>
      </label>
      <label>
        Domain to block:
        <input bind:value={domain} placeholder="e.g., example.com" required>
      </label>
      <button type="submit">Add Rule</button>
    </form>
    
    {#if message}
      <p class="message">{message}</p>
    {/if}
    
    <h2>Blocked Domains</h2>
    {#if rules.length > 0}
      <ul>
        {#each rules as rule}
          <li>
            {rule.domain}
            <button on:click={() => removeRule(rule.id, rule.domain)}>Remove</button>
          </li>
        {/each}
      </ul>
    {:else}
      <p>No rules found for this user.</p>
    {/if}
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
    }

    h1, h2 {
      font-family: 'Playfair Display', serif;
      color: #FFFFFF;
    }

    h1 {
      font-size: 2.5rem;
      margin-bottom: 1rem;
    }

    h2 {
      font-size: 2rem;
      margin-top: 2rem;
    }

    form {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      margin-bottom: 2rem;
    }

    label {
      display: flex;
      flex-direction: column;
      font-size: 1.2rem;
    }

    input {
      padding: 0.8rem;
      margin-top: 0.5rem;
      font-size: 1rem;
      background-color: #FFFFFF;
      color: #000033;
      border: none;
    }

    button {
      padding: 0.8rem 1.5rem;
      background-color: #FFFFFF;
      color: #000033;
      border: none;
      cursor: pointer;
      font-size: 1rem;
      transition: background-color 0.3s ease;
    }

    button:hover {
      background-color: #CCCCCC;
    }

    ul {
      list-style-type: none;
      padding: 0;
    }

    li {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1rem 0;
      border-bottom: 1px solid #FFFFFF;
    }

    .message {
      padding: 1rem;
      background-color: rgba(255, 255, 255, 0.1);
      border-left: 5px solid #FFFFFF;
      margin-bottom: 2rem;
      font-size: 1.2rem;
    }
  </style>