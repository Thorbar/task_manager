<script lang="ts">
  import { onMount } from 'svelte';
  import { writable } from 'svelte/store';
  import ButtonTasks from '$lib/components/ButtonTasks.svelte';
  import { PUBLIC_API_URL, PUBLIC_API_PORT } from '$env/static/public';
  import { t } from 'svelte-i18n';
  import { get } from 'svelte/store';
  import type { Task, User } from '$lib/types';
  import { tasks } from '$lib/stores/tasks'; // Usa esta store centralizada
  import { fetchWithAuth } from '$lib/utils/fetchWithAuth';
  import { goto } from '$app/navigation';


  // Stores to manage tab, tasks list, and user data
  const tab = writable<'tasks' | 'profile'>('tasks');
  const userData = writable<User | null>(null);

  const baseUrl = `${PUBLIC_API_URL}:${PUBLIC_API_PORT}`;

  // Editable user fields
  let editUsername = '';
  let editEmail = '';
  let editName = '';
  let editSurname = '';
  let token: string | null = null;

  // Fetch tasks and user info on component mount
  onMount(async () => {
    token = localStorage.getItem('token');
    if (!token) {
      goto('/');
      return;
    }

    try {
      // Fetch user tasks
      const res = await fetchWithAuth(`${baseUrl}/api/tasks`);
      if (res.ok) {
        const data = await res.json();
        tasks.set(data.tasks);
      } else {
        const errorText = await res.text();
        console.error('Error fetching tasks:', errorText);
      }

      // Fetch user profile data
      const userRes = await fetchWithAuth(`${baseUrl}/api/user`);
      if (userRes.ok) {
        const user = await userRes.json();
        userData.set(user);
        editUsername = user.username;
        editEmail = user.email;
        editName = user.name;
        editSurname = user.surname;
      } else {
        const errorText = await userRes.text();
        console.error('Error fetching user data:', errorText);
      }
    } catch (err) {
      console.error('Error loading data:', err);
    }
  });

  // New task input fields
  let newTitle = '';
  let newDescription = '';

  // Add a new task for the current user
  async function addTask() {
    const res = await fetchWithAuth(`${baseUrl}/api/tasks/create`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: newTitle, description: newDescription, userCreated: editUsername })
    });
    if (res.ok) {
      const newTask = await res.json();
      const currentTasks = get(tasks) ?? [];
      tasks.set([...currentTasks, newTask]);
      newTitle = '';
      newDescription = '';
    } else {
      alert(await res.text());
    }
  }

  // Update an existing task
  async function updateTask(task: Task) {
    try {
      const res = await fetchWithAuth(`${baseUrl}/api/tasks/${task.id}/update`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(task)
      });

      if (!res.ok) {
        const errorText = await res.text();
        alert($t('errorUpdating') + ': ' + errorText);
      } else {
        alert($t('taskUpdated'));
      }
    } catch (error) {
      alert($t('unexpectedUpdateError') + ': ' + error);
    }
  }

  // Delete a task by ID
  async function deleteTask(id: number) {
    const res = await fetchWithAuth(`${baseUrl}/api/tasks/${id}`, {
      method: 'DELETE'
    });
    if (res.ok) {
      tasks.update(ts => ts.filter(t => t.id !== id));
    } else {
      alert($t('errorDeleting') + ': ' + await res.text());
    }
  }

  // Switch between "tasks" and "profile" tabs
  function switchTab(t: 'tasks' | 'profile') {
    tab.set(t);
  }

  // Log out user and redirect to home
  const logout = () => {
    localStorage.removeItem('token');
    location.href = '/';
  };

  // Save user profile changes
  const saveChanges = async () => {
    if (!token) {
      goto('/');
      return;
    }

    const res = await fetchWithAuth(`${baseUrl}/api/user/update`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: editUsername,
        name: editName,
        surname: editSurname
      })
    });

    if (res.ok) {
      alert($t('dataUpdated'));
      userData.set({
        username: editUsername,
        email: editEmail,
        name: editName,
        surname: editSurname,
        created_at: new Date().toISOString()
      });
    } else {
      const err = await res.text();
      alert($t('error') + ': ' + err);
    }
  };
</script>



<section class="flex flex-col items-center">
  <h2 class="text-3xl font-extrabold uppercase tracking-wide text-green-600 drop-shadow-lg mb-10">
    {$tab === 'tasks' ? $t('task') : $t('user')}
  </h2>
  <!-- Navigation buttons -->
  <div class="flex gap-4 mb-4">
    <div class="w-14">
      <ButtonTasks on:click={() => switchTab('tasks')}>
        {$t('task')}
      </ButtonTasks>
    </div>
    <div class="w-16">
      <ButtonTasks on:click={() => switchTab('profile')}>
        {$t('user')}
      </ButtonTasks>
    </div>
    <div class="w-26">
      <ButtonTasks on:click={logout}>
        {$t('logout')}
      </ButtonTasks>
    </div>
  </div>

  <div class="overflow-x-auto border border-gray-300 rounded-t-lg shadow-lg">
    <table class="w-full text-sm text-center bg-gray-400 text-white shadow-lg">
      <thead>
        {#if $tab === 'tasks'}
          <tr>
            <th class="px-4 py-2 border">{$t('title')}</th>
            <th class="px-4 py-2 border">{$t('description')}</th>
            <th class="px-4 py-2 border">{$t('user')}</th>
            <th class="px-4 py-2 border">{$t('created_at')}</th>
            <th class="px-4 py-2 border">{$t('status')}</th>
            <th class="px-4 py-2 border">{$t('options')}</th>
          </tr>
        {:else}
          <tr>
            <th class="px-4 py-2 border">{$t('user')}</th>
            <th class="px-4 py-2 border">{$t('name')}</th>
            <th class="px-4 py-2 border">{$t('surname')}</th>
            <th class="px-4 py-2 border">{$t('email')}</th>
            <th class="px-4 py-2 border">{$t('created_at')}</th>
            <th class="px-4 py-2 border"></th>
          </tr>
        {/if}
      </thead>
      <tbody>
        {#if $tab === 'tasks'}
          {#each $tasks as task}
            <tr class="bg-yellow-100 text-black border-b hover:bg-yellow-200">
              <td class="px-4 py-2 border"><input class="bg-transparent w-full text-center" bind:value={task.title} /></td>
              <td class="px-4 py-2"><input class="bg-transparent w-full text-center" bind:value={task.description} /></td>
              <td class="px-4 py-2 border text-gray-600">{task.username}</td>
              <td class="px-4 py-2 border">{new Date(task.created_at).toLocaleDateString()}</td>
              <td class="px-4 py-2 border font-semibold {task.status === 'completed' ? 'text-green-600' : 'text-red-600'}">
                <select bind:value={task.status} class ="bg-transparent">
                  <option class="text-red-600" value="pending">{$t('taskPending')}</option>
                  <option class="text-green-600" value="completed">{$t('taskCompleted')}</option>
                </select>
              </td>
              <td class="flex gap-2 mt-1 font-bold hover:text-blue-500 transition duration-300 transform">
                <button on:click={() => updateTask(task)} class="bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600">
                  {$t('save')}
                </button>
                <button on:click={() => deleteTask(task.id)} class="bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600">
                  {$t('delete')}
                </button>
              </td>
            </tr>
          {/each}

          <!-- New task input row -->
          <tr class="bg-cyan-100 text-black border-b hover:bg-yellow-200">
            <td><input bind:value={newTitle} placeholder={$t('title')} class="bg-transparent text-center " /></td>
            <td><input bind:value={newDescription} placeholder={$t('description')} class="bg-transparent text-center " /></td>
            <td class="px-4 py-2 border text-gray-600">{editUsername}</td>
            <td class="px-4 py-2 border"></td>
            <td class="px-4 py-2 border">{$t('taskPending')}</td>
            <td class="border font-bold hover:text-blue-500 transition duration-300 transform" colspan="1">
              <button on:click={addTask}>
                {$t('add')}
              </button>
            </td>
          </tr>
        {:else if $tab === 'profile' && $userData}
          <!-- User profile edit row -->
          <tr class="bg-yellow-100 text-black border-b hover:bg-yellow-200">
            <td class="px-4 py-2 border">{editUsername}</td>
            <td class="px-4 py-2 border">
              <input type="text" bind:value={editName} placeholder={$t('name')} class = "bg-transparent"/>
            </td>
            <td class="px-4 py-2 border">
              <input type="text" bind:value={editSurname} placeholder={$t('surname')} class = "bg-transparent" />
            </td>
            <td class="px-4 py-2 border">{editEmail}</td>
            <td class="px-4 py-2 border">{new Date($userData.created_at).toLocaleDateString()}</td>
            <td class="px-4 py-2 font-bold hover:text-blue-500 transition duration-300 transform">
              <button on:click={saveChanges}>
                {$t('save')}
              </button>
            </td>
          </tr>
        {/if}
      </tbody>
    </table>
  </div>
</section>
