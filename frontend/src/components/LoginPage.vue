<template>
    <div class="min-h-screen bg-base-100 flex items-center justify-center p-4">
        <div class="w-full max-w-sm">
            <h1 class="text-2xl font-semibold text-center mb-8">Punchcard</h1>
            <form @submit.prevent="login" class="rounded-lg border border-base-content/10 bg-base-200 p-6 space-y-4">
                <div>
                    <label class="text-xs text-base-content/50 block mb-1.5">Username</label>
                    <input
                        v-model="username"
                        type="text"
                        autocomplete="username"
                        class="input input-sm input-bordered w-full bg-base-300 border-base-content/10 focus:border-base-content/20 focus:outline-none text-sm"
                    />
                </div>
                <div>
                    <label class="text-xs text-base-content/50 block mb-1.5">Password</label>
                    <input
                        v-model="password"
                        type="password"
                        autocomplete="current-password"
                        class="input input-sm input-bordered w-full bg-base-300 border-base-content/10 focus:border-base-content/20 focus:outline-none text-sm"
                    />
                </div>
                <p v-if="error" class="text-xs text-error">{{ error }}</p>
                <button
                    type="submit"
                    :disabled="loading"
                    class="w-full text-sm py-2 rounded-lg bg-primary text-primary-content hover:opacity-90 transition-opacity cursor-pointer disabled:opacity-50"
                >
                    {{ loading ? "Signing in..." : "Sign in" }}
                </button>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const emit = defineEmits<{ authenticated: [] }>();

const username = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);

async function login() {
    error.value = "";
    loading.value = true;

    try {
        const res = await fetch("/web/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                username: username.value,
                password: password.value,
            }),
        });

        if (res.ok) {
            emit("authenticated");
        } else {
            error.value = "Invalid username or password";
        }
    } catch {
        error.value = "Connection error";
    } finally {
        loading.value = false;
    }
}
</script>
