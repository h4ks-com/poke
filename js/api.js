// Viridian City Bank API Client
class BankAPI {
    constructor() {
        this.baseURL = '/api'; // Same origin since frontend and backend are served together
        this.token = localStorage.getItem('authToken');
    }

    // Get authorization headers
    getHeaders() {
        const headers = {
            'Content-Type': 'application/json'
        };
        
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        
        return headers;
    }

    // Handle API responses
    async handleResponse(response) {
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || `HTTP error! status: ${response.status}`);
        }
        
        return data;
    }

    // Set authentication token
    setToken(token) {
        this.token = token;
        if (token) {
            localStorage.setItem('authToken', token);
        } else {
            localStorage.removeItem('authToken');
        }
    }

    // Authentication endpoints
    async register(username, email, password, confirmPassword) {
        console.log(`[API] POST ${this.baseURL}/register`);
        
        const response = await fetch(`${this.baseURL}/register`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({
                username,
                email,
                password,
                confirmPassword
            })
        });

        const data = await this.handleResponse(response);
        return data;
    }

    async login(username, password) {
        console.log(`[API] POST ${this.baseURL}/login`);
        
        const response = await fetch(`${this.baseURL}/login`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({
                username,
                password
            })
        });

        const data = await this.handleResponse(response);
        
        if (data.success && data.token) {
            this.setToken(data.token);
            localStorage.setItem('currentUser', JSON.stringify(data.user));
        }
        
        return data;
    }

    async changePassword(currentPassword, newPassword, confirmNewPassword) {
        console.log(`[API] POST ${this.baseURL}/change-password`);
        
        const response = await fetch(`${this.baseURL}/change-password`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({
                currentPassword,
                newPassword,
                confirmNewPassword
            })
        });

        return await this.handleResponse(response);
    }

    async logout() {
        console.log(`[API] Logging out`);
        this.setToken(null);
        localStorage.removeItem('currentUser');
        return { success: true, message: 'Logged out successfully' };
    }

    // Banking endpoints
    async getBalance() {
        console.log(`[API] GET ${this.baseURL}/balance`);
        
        const response = await fetch(`${this.baseURL}/balance`, {
            method: 'GET',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    async getAccountInfo() {
        console.log(`[API] GET ${this.baseURL}/account`);
        
        const response = await fetch(`${this.baseURL}/account`, {
            method: 'GET',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    async getTransactions(limit = 50) {
        console.log(`[API] GET ${this.baseURL}/transactions?limit=${limit}`);
        
        const response = await fetch(`${this.baseURL}/transactions?limit=${limit}`, {
            method: 'GET',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    async transfer(to, amount, description = '') {
        console.log(`[API] POST ${this.baseURL}/transfer`);
        
        const response = await fetch(`${this.baseURL}/transfer`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({
                to,
                amount: parseFloat(amount),
                description
            })
        });

        return await this.handleResponse(response);
    }

    // Payment request endpoints
    async createPaymentRequest(to, amount, reason, message = '') {
        console.log(`[API] POST ${this.baseURL}/payment-requests`);
        
        const response = await fetch(`${this.baseURL}/payment-requests`, {
            method: 'POST',
            headers: this.getHeaders(),
            body: JSON.stringify({
                to,
                amount: parseFloat(amount),
                reason,
                message
            })
        });

        return await this.handleResponse(response);
    }

    async getPaymentRequests() {
        console.log(`[API] GET ${this.baseURL}/payment-requests`);
        
        const response = await fetch(`${this.baseURL}/payment-requests`, {
            method: 'GET',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    async handlePaymentRequest(requestId, action) {
        console.log(`[API] PUT ${this.baseURL}/payment-requests/${requestId}`);
        
        const response = await fetch(`${this.baseURL}/payment-requests/${requestId}`, {
            method: 'PUT',
            headers: this.getHeaders(),
            body: JSON.stringify({
                action // 'approve' or 'reject'
            })
        });

        return await this.handleResponse(response);
    }

    // Card endpoints
    async getCard() {
        console.log(`[API] GET ${this.baseURL}/card`);
        
        const response = await fetch(`${this.baseURL}/card`, {
            method: 'GET',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    async refreshCard() {
        console.log(`[API] POST ${this.baseURL}/card/refresh`);
        
        const response = await fetch(`${this.baseURL}/card/refresh`, {
            method: 'POST',
            headers: this.getHeaders()
        });

        return await this.handleResponse(response);
    }

    // Health check endpoint
    async healthCheck() {
        console.log(`[API] GET ${this.baseURL.replace('/api', '')}/health`);
        
        const response = await fetch(`${this.baseURL.replace('/api', '')}/health`, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' }
        });

        return await this.handleResponse(response);
    }
}

// Create global API instance
const api = new BankAPI();

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
    module.exports = BankAPI;
}
