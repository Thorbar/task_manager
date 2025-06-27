<script lang="ts">
  import { onMount } from 'svelte';
  import { tasks, errorMsg } from '$lib/stores/tasks';
  import { t } from 'svelte-i18n';

  /**
   * On component mount, fetch the tasks of the authenticated user.
   * If the user is not logged in (no token), set an appropriate error message.
   * Handles success and failure responses, setting tasks or error messages accordingly.
   * Network errors are caught and reported as well.
   */
  onMount(async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      // User not authenticated, set error message
      errorMsg.set($t('notLoggedIn'));
      return;
    }

    try {
      const res = await fetch('/api/tasks', {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (res.ok) {
        const data = await res.json();
        tasks.set(data.tasks);
        errorMsg.set(''); // Clear error message on success
      } else {
        // Token invalid or expired
        errorMsg.set($t('invalidToken'));
      }
    } catch (e) {
      // Network or unexpected error
      errorMsg.set($t('networkError'));
    }
  });
</script>
