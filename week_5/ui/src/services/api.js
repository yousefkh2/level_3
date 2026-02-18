// JWT Management Suite (well-structured)

import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL

export const getToken = () => {
    return localStorage.getItem('token')
}
export const setToken = (token) => {
    localStorage.setItem('token', token)
}

export const clearToken = () => {
    localStorage.removeItem('token')
}


export const apiClient = axios.create({
    baseURL: API_BASE_URL
})


// add token to every request
apiClient.interceptors.request.use((config) => {
    const token = getToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
})

export default apiClient;