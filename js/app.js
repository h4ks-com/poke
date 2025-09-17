// Main Banking Application
class BankingApp {
    constructor() {
        this.api = new BankAPI();
        this.currentUser = null;
        this.isLoggedIn = false;
        this.cardRefreshSeed = 0; // For generating new card numbers
        
        this.initializeApp();
        this.bindEvents();
    }

    initializeApp() {
        // Check if user is logged in
        if (localStorage.getItem('authToken')) {
            // Restore user data from localStorage
            const storedUser = localStorage.getItem('currentUser');
            if (storedUser) {
                this.currentUser = JSON.parse(storedUser);
                this.isLoggedIn = true;
            }
            this.showDashboard();
        } else {
            this.showLogin();
        }
    }

    bindEvents() {
        // Login form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) {
            loginForm.addEventListener('submit', (e) => this.handleLogin(e));
        }

        // Logout button
        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', () => this.handleLogout());
        }

        // User menu dropdown
        const userMenuBtn = document.getElementById('userMenuBtn');
        const userDropdown = document.getElementById('userDropdown');
        if (userMenuBtn && userDropdown) {
            userMenuBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                userDropdown.classList.toggle('show');
            });
            
            // Close dropdown when clicking outside
            document.addEventListener('click', () => {
                userDropdown.classList.remove('show');
            });
        }

        // Registration and password change buttons
        const showRegisterBtn = document.getElementById('showRegisterBtn');
        const changePasswordBtn = document.getElementById('changePasswordBtn');
        
        if (showRegisterBtn) {
            showRegisterBtn.addEventListener('click', () => this.openRegistrationModal());
        }
        if (changePasswordBtn) {
            changePasswordBtn.addEventListener('click', () => this.openPasswordModal());
        }

        // Action buttons
        const transferBtn = document.getElementById('transferBtn');
        const requestBtn = document.getElementById('requestBtn');
        const viewCardBtn = document.getElementById('viewCardBtn');

        if (transferBtn) transferBtn.addEventListener('click', () => this.openTransferModal());
        if (requestBtn) requestBtn.addEventListener('click', () => this.openRequestModal());
        if (viewCardBtn) viewCardBtn.addEventListener('click', () => this.openCardModal());

        // Modal events
        this.bindModalEvents();

        // Payment requests tabs
        this.bindPaymentRequestTabs();

        // Form submissions
        this.bindFormEvents();
    }

    bindModalEvents() {
        // Request modal
        const requestModal = document.getElementById('requestModal');
        const closeRequestModal = document.getElementById('closeRequestModal');
        const cancelRequest = document.getElementById('cancelRequest');

        if (closeRequestModal) closeRequestModal.addEventListener('click', () => this.closeModal('requestModal'));
        if (cancelRequest) cancelRequest.addEventListener('click', () => this.closeModal('requestModal'));

        // Transfer modal
        const transferModal = document.getElementById('transferModal');
        const closeTransferModal = document.getElementById('closeTransferModal');
        const cancelTransfer = document.getElementById('cancelTransfer');

        if (closeTransferModal) closeTransferModal.addEventListener('click', () => this.closeModal('transferModal'));
        if (cancelTransfer) cancelTransfer.addEventListener('click', () => this.closeModal('transferModal'));

        // Registration modal
        const registrationModal = document.getElementById('registrationModal');
        const closeRegistrationModal = document.getElementById('closeRegistrationModal');
        const cancelRegistration = document.getElementById('cancelRegistration');

        if (closeRegistrationModal) closeRegistrationModal.addEventListener('click', () => this.closeModal('registrationModal'));
        if (cancelRegistration) cancelRegistration.addEventListener('click', () => this.closeModal('registrationModal'));

        // Password change modal
        const passwordModal = document.getElementById('passwordModal');
        const closePasswordModal = document.getElementById('closePasswordModal');
        const cancelPasswordChange = document.getElementById('cancelPasswordChange');

        if (closePasswordModal) closePasswordModal.addEventListener('click', () => this.closeModal('passwordModal'));
        if (cancelPasswordChange) cancelPasswordChange.addEventListener('click', () => this.closeModal('passwordModal'));

        // Card modal
        const cardModal = document.getElementById('cardModal');
        const closeCardModal = document.getElementById('closeCardModal');

        if (closeCardModal) closeCardModal.addEventListener('click', () => this.closeModal('cardModal'));

        // Wallet buttons
        const addToAppleWallet = document.getElementById('addToAppleWallet');
        const addToGooglePay = document.getElementById('addToGooglePay');
        const refreshCardBtn = document.getElementById('refreshCardBtn');

        if (addToAppleWallet) addToAppleWallet.addEventListener('click', () => this.addToAppleWallet());
        if (addToGooglePay) addToGooglePay.addEventListener('click', () => this.addToGooglePay());
        if (refreshCardBtn) refreshCardBtn.addEventListener('click', () => this.refreshCard());

        // Close modals when clicking outside
        const modals = ['requestModal', 'transferModal', 'registrationModal', 'passwordModal', 'cardModal'];
        modals.forEach(modalId => {
            const modal = document.getElementById(modalId);
            if (modal) {
                modal.addEventListener('click', (e) => {
                    if (e.target === modal) {
                        this.closeModal(modalId);
                    }
                });
            }
        });
    }

    bindPaymentRequestTabs() {
        const incomingTab = document.getElementById('incomingTab');
        const outgoingTab = document.getElementById('outgoingTab');
        const incomingRequests = document.getElementById('incomingRequests');
        const outgoingRequests = document.getElementById('outgoingRequests');

        if (incomingTab) {
            incomingTab.addEventListener('click', () => {
                incomingTab.classList.add('active');
                outgoingTab.classList.remove('active');
                incomingRequests.style.display = 'block';
                outgoingRequests.style.display = 'none';
            });
        }

        if (outgoingTab) {
            outgoingTab.addEventListener('click', () => {
                outgoingTab.classList.add('active');
                incomingTab.classList.remove('active');
                outgoingRequests.style.display = 'block';
                incomingRequests.style.display = 'none';
            });
        }
    }

    bindFormEvents() {
        // Request form
        const requestForm = document.getElementById('requestForm');
        if (requestForm) {
            requestForm.addEventListener('submit', (e) => this.handleRequest(e));
        }

        // Transfer form
        const transferForm = document.getElementById('transferForm');
        if (transferForm) {
            transferForm.addEventListener('submit', (e) => this.handleTransfer(e));
        }

        // Registration form
        const registrationForm = document.getElementById('registrationForm');
        if (registrationForm) {
            registrationForm.addEventListener('submit', (e) => this.handleRegistration(e));
        }

        // Password change form
        const passwordChangeForm = document.getElementById('passwordChangeForm');
        if (passwordChangeForm) {
            passwordChangeForm.addEventListener('submit', (e) => this.handlePasswordChange(e));
        }
    }

    async handleLogin(e) {
        e.preventDefault();
        
        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        const loginBtn = e.target.querySelector('button[type="submit"]');

        // Show loading state
        const originalText = loginBtn.innerHTML;
        loginBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Signing In...';
        loginBtn.disabled = true;

        try {
            const response = await this.api.login(username, password);
            
            this.currentUser = response.user;
            this.isLoggedIn = true;
            
            // Save user session (for demo)
            localStorage.setItem('currentUser', JSON.stringify(this.currentUser));
            
            this.showMessage('success', 'Login Successful', `Welcome back, ${this.currentUser.username}!`);
            this.showDashboard();
            
        } catch (error) {
            this.showMessage('error', 'Login Failed', 'Invalid username or password');
        } finally {
            loginBtn.innerHTML = originalText;
            loginBtn.disabled = false;
        }
    }

    async handleLogout() {
        try {
            await this.api.logout();
            
            this.currentUser = null;
            this.isLoggedIn = false;
            
            localStorage.removeItem('currentUser');
            
            this.showMessage('info', 'Logged Out', 'You have been successfully logged out');
            this.showLogin();
            
        } catch (error) {
            this.showMessage('error', 'Logout Error', 'Failed to logout properly');
        }
    }

    async handleRequest(e) {
        e.preventDefault();
        
        const toUsername = document.getElementById('requestFromAccount').value;
        const amount = parseFloat(document.getElementById('requestAmount').value);
        const reason = document.getElementById('requestReason').value;
        const message = document.getElementById('requestMessage').value;
        
        const submitBtn = e.target.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Sending...';
        submitBtn.disabled = true;

        try {
            const response = await this.api.createPaymentRequest(toUsername, amount, reason, message);
            
            this.showMessage('success', 'Request Sent', `Payment request for ₽${amount.toFixed(2)} sent to ${toUsername}`);
            this.closeModal('requestModal');
            this.refreshDashboard();
            
        } catch (error) {
            console.error('Payment request error:', error);
            this.showMessage('error', 'Request Failed', error.message || 'Failed to send payment request');
        } finally {
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    }

    async handleTransfer(e) {
        e.preventDefault();
        
        const recipientAccount = document.getElementById('recipientAccount').value;
        const amount = parseFloat(document.getElementById('transferAmount').value);
        const memo = document.getElementById('transferMemo').value;
        
        const submitBtn = e.target.querySelector('button[type="submit"]');
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Processing...';
        submitBtn.disabled = true;

        try {
            const response = await this.api.transfer(recipientAccount, amount, memo);
            
            this.showMessage('success', 'Transfer Successful', `$${amount.toFixed(2)} sent to ${recipientAccount}`);
            this.closeModal('transferModal');
            this.refreshDashboard();
            
        } catch (error) {
            this.showMessage('error', 'Transfer Failed', 'Insufficient funds or invalid recipient');
        } finally {
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    }

    showLogin() {
        document.getElementById('loginScreen').style.display = 'flex';
        document.getElementById('bankingDashboard').style.display = 'none';
        
        // Clear form
        const loginForm = document.getElementById('loginForm');
        if (loginForm) loginForm.reset();
    }

    bindLoginScreenEvents() {
        // Registration button
        const showRegisterBtn = document.getElementById('showRegisterBtn');
        console.log('bindLoginScreenEvents - showRegisterBtn found:', showRegisterBtn);
        
        if (showRegisterBtn) {
            // Remove any existing listeners
            showRegisterBtn.replaceWith(showRegisterBtn.cloneNode(true));
            const newShowRegisterBtn = document.getElementById('showRegisterBtn');
            
            console.log('Adding fresh click listener to showRegisterBtn');
            newShowRegisterBtn.addEventListener('click', (e) => {
                e.preventDefault();
                console.log('showRegisterBtn clicked!');
                this.openRegistrationModal();
            });
        }
    }

    showDashboard() {
        document.getElementById('loginScreen').style.display = 'none';
        document.getElementById('bankingDashboard').style.display = 'block';
        
        // Update user info
        const currentUserEl = document.getElementById('currentUser');
        if (currentUserEl && this.currentUser) {
            currentUserEl.textContent = this.currentUser.username;
        }
        
        // Update account number display
        const accountNumberEl = document.getElementById('accountNumberDisplay');
        const toggleBtn = document.getElementById('toggleAccountBtn');
        const copyBtn = document.getElementById('copyAccountBtn');
        
        if (accountNumberEl && this.currentUser) {
            const accountNumber = this.currentUser.account_number;
            if (accountNumber) {
                // Store the full account number for toggle functionality
                this.fullAccountNumber = accountNumber;
                this.isAccountNumberRevealed = false;
                
                // Format account number with partial masking (show last 4 digits)
                const maskedAccountNumber = this.formatAccountNumber(accountNumber);
                accountNumberEl.textContent = `Account: ${maskedAccountNumber}`;
                
                // Show the toggle and copy buttons
                if (toggleBtn) toggleBtn.style.display = 'flex';
                if (copyBtn) copyBtn.style.display = 'flex';
                
                // Bind event listeners
                this.bindAccountNumberEvents();
            } else {
                accountNumberEl.textContent = 'Account: Not Available';
            }
        }
        
        this.refreshDashboard();
    }

    formatDate(dateString) {
        if (!dateString) return 'Unknown Date';
        
        try {
            // Try parsing the date string
            const date = new Date(dateString);
            
            // Check if the date is valid
            if (isNaN(date.getTime())) {
                // If invalid, try parsing as ISO string or handle SQLite format
                const sqliteDate = dateString.replace(' ', 'T');
                const retryDate = new Date(sqliteDate);
                
                if (isNaN(retryDate.getTime())) {
                    console.warn('Could not parse date:', dateString);
                    return 'Invalid Date';
                }
                
                return retryDate.toLocaleDateString() + ' ' + retryDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
            }
            
            return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
        } catch (error) {
            console.error('Date formatting error:', error, dateString);
            return 'Invalid Date';
        }
    }

    formatAccountNumber(accountNumber) {
        if (!accountNumber) return '****';
        
        // Convert to string if it's a number
        const accountStr = accountNumber.toString();
        
        if (accountStr.length <= 4) {
            // If account number is 4 digits or less, show it all
            return accountStr;
        } else {
            // Show first 2 and last 4 digits with asterisks in between
            const firstTwo = accountStr.substring(0, 2);
            const lastFour = accountStr.substring(accountStr.length - 4);
            const middleAsterisks = '*'.repeat(Math.max(0, accountStr.length - 6));
            return `${firstTwo}${middleAsterisks}${lastFour}`;
        }
    }

    bindAccountNumberEvents() {
        const toggleBtn = document.getElementById('toggleAccountBtn');
        const copyBtn = document.getElementById('copyAccountBtn');
        
        // Remove existing listeners to prevent duplicates
        if (toggleBtn) {
            toggleBtn.replaceWith(toggleBtn.cloneNode(true));
            const newToggleBtn = document.getElementById('toggleAccountBtn');
            newToggleBtn.addEventListener('click', () => this.toggleAccountNumber());
        }
        
        if (copyBtn) {
            copyBtn.replaceWith(copyBtn.cloneNode(true));
            const newCopyBtn = document.getElementById('copyAccountBtn');
            newCopyBtn.addEventListener('click', () => this.copyAccountNumber());
        }
    }

    toggleAccountNumber() {
        const accountNumberEl = document.getElementById('accountNumberDisplay');
        const toggleBtn = document.getElementById('toggleAccountBtn');
        const eyeIcon = toggleBtn.querySelector('i');
        
        if (!this.fullAccountNumber || !accountNumberEl) return;
        
        this.isAccountNumberRevealed = !this.isAccountNumberRevealed;
        
        if (this.isAccountNumberRevealed) {
            // Show full account number
            accountNumberEl.textContent = `Account: ${this.fullAccountNumber}`;
            eyeIcon.className = 'fas fa-eye-slash';
            toggleBtn.title = 'Hide account number';
        } else {
            // Show masked account number
            const maskedAccountNumber = this.formatAccountNumber(this.fullAccountNumber);
            accountNumberEl.textContent = `Account: ${maskedAccountNumber}`;
            eyeIcon.className = 'fas fa-eye';
            toggleBtn.title = 'Show full account number';
        }
    }

    async copyAccountNumber() {
        const copyBtn = document.getElementById('copyAccountBtn');
        const copyIcon = copyBtn.querySelector('i');
        
        if (!this.fullAccountNumber) return;
        
        try {
            await navigator.clipboard.writeText(this.fullAccountNumber);
            
            // Visual feedback
            copyIcon.className = 'fas fa-check';
            copyBtn.classList.add('copied');
            copyBtn.title = 'Copied!';
            
            // Reset after 2 seconds
            setTimeout(() => {
                copyIcon.className = 'fas fa-copy';
                copyBtn.classList.remove('copied');
                copyBtn.title = 'Copy account number';
            }, 2000);
            
        } catch (error) {
            console.error('Failed to copy account number:', error);
            // Fallback for older browsers
            this.fallbackCopyAccountNumber();
        }
    }

    fallbackCopyAccountNumber() {
        // Fallback method for browsers that don't support clipboard API
        const textArea = document.createElement('textarea');
        textArea.value = this.fullAccountNumber;
        document.body.appendChild(textArea);
        textArea.select();
        
        try {
            document.execCommand('copy');
            this.showMessage('success', 'Copied', 'Account number copied to clipboard');
        } catch (error) {
            this.showMessage('error', 'Copy Failed', 'Could not copy account number');
        }
        
        document.body.removeChild(textArea);
    }

    async refreshDashboard() {
        await this.updateBalance();
        await this.loadTransactions();
        await this.loadPaymentRequests();
    }

    async updateBalance() {
        try {
            const response = await this.api.getBalance();
            const balanceEl = document.getElementById('accountBalance');
            if (balanceEl) {
                balanceEl.textContent = response.balance.toFixed(2);
            }
        } catch (error) {
            console.error('Failed to update balance:', error);
            const balanceEl = document.getElementById('accountBalance');
            if (balanceEl) {
                balanceEl.textContent = '0.00';
            }
        }
    }

    async loadTransactions() {
        try {
            const response = await this.api.getTransactions(10);
            const transactionsContainer = document.getElementById('transactionsList');
            
            if (transactionsContainer) {
                if (response.transactions.length === 0) {
                    transactionsContainer.innerHTML = '<p style="text-align: center; color: #666; padding: 20px;">No transactions yet</p>';
                } else {
                    transactionsContainer.innerHTML = response.transactions.map(tx => this.createTransactionHTML(tx)).join('');
                }
            }
        } catch (error) {
            console.error('Failed to load transactions:', error);
        }
    }

    async loadPaymentRequests() {
        try {
            const [incomingResponse, outgoingResponse] = await Promise.all([
                this.api.getPaymentRequests('incoming'),
                this.api.getPaymentRequests('outgoing')
            ]);
            
            // Update counters
            const incomingCount = document.getElementById('incomingCount');
            const outgoingCount = document.getElementById('outgoingCount');
            
            if (incomingCount) {
                incomingCount.textContent = incomingResponse.requests.filter(req => req.status === 'pending').length;
            }
            if (outgoingCount) {
                outgoingCount.textContent = outgoingResponse.requests.filter(req => req.status === 'pending').length;
            }
            
            // Update request lists
            const incomingContainer = document.getElementById('incomingRequests');
            const outgoingContainer = document.getElementById('outgoingRequests');
            
            if (incomingContainer) {
                if (incomingResponse.requests.length === 0) {
                    incomingContainer.innerHTML = '<div class="empty-state"><i class="fas fa-inbox"></i><p>No payment requests received</p></div>';
                } else {
                    incomingContainer.innerHTML = incomingResponse.requests.map(req => this.createPaymentRequestHTML(req, 'incoming')).join('');
                }
            }
            
            if (outgoingContainer) {
                if (outgoingResponse.requests.length === 0) {
                    outgoingContainer.innerHTML = '<div class="empty-state"><i class="fas fa-paper-plane"></i><p>No payment requests sent</p></div>';
                } else {
                    outgoingContainer.innerHTML = outgoingResponse.requests.map(req => this.createPaymentRequestHTML(req, 'outgoing')).join('');
                }
            }
            
            // Bind action buttons
            this.bindPaymentRequestActions();
            
        } catch (error) {
            console.error('Failed to load payment requests:', error);
        }
    }

    createPaymentRequestHTML(request, type) {
        const isIncoming = type === 'incoming';
        const otherAccount = isIncoming ? request.from_username || request.fromUsername : request.to_username || request.toUsername;
        const formattedDate = this.formatDate(request.created_at || request.createdAt);
        
        const statusClass = request.status;
        const statusText = request.status.charAt(0).toUpperCase() + request.status.slice(1);
        
        let actionsHTML = '';
        if (request.status === 'pending') {
            if (isIncoming) {
                actionsHTML = `
                    <div class="request-actions">
                        <button type="button" class="accept-btn" data-request-id="${request.id}" data-action="accept">
                            <i class="fas fa-check"></i> Accept
                        </button>
                        <button type="button" class="reject-btn" data-request-id="${request.id}" data-action="reject">
                            <i class="fas fa-times"></i> Reject
                        </button>
                    </div>
                `;
            } else {
                actionsHTML = `
                    <div class="request-actions">
                        <button type="button" class="cancel-request-btn" data-request-id="${request.id}" data-action="cancel">
                            <i class="fas fa-ban"></i> Cancel
                        </button>
                    </div>
                `;
            }
        }
        
        return `
            <div class="payment-request-item ${type}">
                <div class="payment-request-header">
                    <div class="request-info">
                        <h4>${isIncoming ? 'Request from' : 'Request to'} ${otherAccount}</h4>
                        <p><strong>Reason:</strong> ${request.reason}</p>
                        <p><strong>Date:</strong> ${formattedDate}</p>
                    </div>
                    <div>
                        <div class="request-amount">$${request.amount.toFixed(2)}</div>
                        <div class="request-status ${statusClass}">${statusText}</div>
                    </div>
                </div>
                ${request.message ? `<div class="request-message">${request.message}</div>` : ''}
                ${actionsHTML}
            </div>
        `;
    }

    bindPaymentRequestActions() {
        // Accept/Reject buttons
        document.querySelectorAll('.accept-btn, .reject-btn').forEach(btn => {
            btn.addEventListener('click', async (e) => {
                try {
                    e.preventDefault();
                    e.stopPropagation();
                    const requestId = e.currentTarget.dataset.requestId;
                    const action = e.currentTarget.dataset.action;
                    await this.handlePaymentRequestResponse(requestId, action, e.currentTarget);
                } catch (error) {
                    console.error('Error in accept/reject button handler:', error);
                    e.preventDefault();
                    e.stopPropagation();
                }
            });
        });
        
        // Cancel buttons
        document.querySelectorAll('.cancel-request-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                e.stopImmediatePropagation();
                
                const requestId = e.currentTarget.dataset.requestId;
                this.handlePaymentRequestCancel(requestId, e.currentTarget);
            });
        });
    }

    async refreshActivePaymentRequestTab() {
        // Determine which tab is active
        const incomingTab = document.getElementById('incomingTab');
        const outgoingTab = document.getElementById('outgoingTab');
        
        const isIncomingActive = incomingTab && incomingTab.classList.contains('active');
        const isOutgoingActive = outgoingTab && outgoingTab.classList.contains('active');
        
        try {
            const response = await this.api.getPaymentRequests();
            
            // Only update the active tab
            if (isIncomingActive) {
                const incomingRequests = response.incoming || [];
                this.updatePaymentRequestsDisplay(incomingRequests, 'incoming');
                
                // Update counter
                const incomingCount = incomingRequests.filter(req => req.status === 'pending').length;
                document.getElementById('incomingCount').textContent = incomingCount.toString();
            } else if (isOutgoingActive) {
                const outgoingRequests = response.outgoing || [];
                this.updatePaymentRequestsDisplay(outgoingRequests, 'outgoing');
                
                // Update counter
                const outgoingCount = outgoingRequests.filter(req => req.status === 'pending').length;
                document.getElementById('outgoingCount').textContent = outgoingCount.toString();
            }
            
        } catch (error) {
            console.error('Failed to refresh active payment request tab:', error);
        }
    }

    async loadPaymentRequests() {
        try {
            const response = await this.api.getPaymentRequests();
            
            // The backend returns { incoming: [...], outgoing: [...] }
            const incomingRequests = response.incoming || [];
            const outgoingRequests = response.outgoing || [];

            this.updatePaymentRequestsDisplay(incomingRequests, 'incoming');
            this.updatePaymentRequestsDisplay(outgoingRequests, 'outgoing');

            // Update counters
            const incomingCount = incomingRequests.filter(req => req.status === 'pending').length;
            const outgoingCount = outgoingRequests.filter(req => req.status === 'pending').length;

            document.getElementById('incomingCount').textContent = incomingCount;
            document.getElementById('outgoingCount').textContent = outgoingCount;

        } catch (error) {
            console.error('Failed to load payment requests:', error);
            // Set empty arrays as fallback
            this.updatePaymentRequestsDisplay([], 'incoming');
            this.updatePaymentRequestsDisplay([], 'outgoing');
            document.getElementById('incomingCount').textContent = '0';
            document.getElementById('outgoingCount').textContent = '0';
        }
    }

    updatePaymentRequestsDisplay(requests, type) {
        const container = document.getElementById(type === 'incoming' ? 'incomingRequests' : 'outgoingRequests');
        
        if (!container) return;

        // Ensure requests is an array
        if (!Array.isArray(requests)) {
            console.warn(`Expected requests to be an array, got:`, requests);
            requests = [];
        }

        if (requests.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-inbox"></i>
                    <p>No ${type} payment requests</p>
                </div>
            `;
        } else {
            container.innerHTML = requests.map(request => this.createPaymentRequestHTML(request, type)).join('');
        }

        // Bind action buttons
        this.bindPaymentRequestActions(container);
    }

    createPaymentRequestHTML(request, type) {
        const formattedDate = this.formatDate(request.created_at || request.createdAt);
        
        const isIncoming = type === 'incoming';
        const otherParty = isIncoming ? request.from_username || request.fromUsername : request.to_username || request.toUsername;
        const canRespond = isIncoming && request.status === 'pending';
        const canCancel = !isIncoming && request.status === 'pending';

        return `
            <div class="payment-request-item ${type}">
                <div class="payment-request-header">
                    <div class="request-info">
                        <h4>${isIncoming ? 'Request from' : 'Request to'} ${otherParty}</h4>
                        <p>${request.reason}</p>
                        <p>Created: ${formattedDate}</p>
                    </div>
                    <div style="text-align: right;">
                        <div class="request-amount">$${request.amount.toFixed(2)}</div>
                        <span class="request-status ${request.status}">${request.status}</span>
                    </div>
                </div>
                
                ${request.message ? `
                    <div class="request-message">
                        "${request.message}"
                    </div>
                ` : ''}
                
                <div class="request-actions">
                    ${canRespond ? `
                        <button type="button" class="accept-btn" data-request-id="${request.id}" data-action="accept">
                            <i class="fas fa-check"></i> Accept
                        </button>
                        <button type="button" class="reject-btn" data-request-id="${request.id}" data-action="reject">
                            <i class="fas fa-times"></i> Reject
                        </button>
                    ` : ''}
                    ${canCancel ? `
                        <button type="button" class="cancel-request-btn" data-request-id="${request.id}" data-action="cancel">
                            <i class="fas fa-ban"></i> Cancel
                        </button>
                    ` : ''}
                </div>
            </div>
        `;
    }

    bindPaymentRequestActions(container) {
        // Accept buttons
        container.querySelectorAll('.accept-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.handlePaymentRequestResponse(e.target.dataset.requestId, 'approve', e.target);
            });
        });

        // Reject buttons
        container.querySelectorAll('.reject-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.handlePaymentRequestResponse(e.target.dataset.requestId, 'reject', e.target);
            });
        });

        // Cancel buttons
        container.querySelectorAll('.cancel-request-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.handlePaymentRequestCancel(e.target.dataset.requestId, e.target);
            });
        });
    }

    async handlePaymentRequestResponse(requestId, action, button) {
        const originalText = button.innerHTML;
        button.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';
        button.disabled = true;

        try {
            const response = await this.api.handlePaymentRequest(requestId, action);
            console.log('Server response:', response); // This proves the server processed it
            
            const actionText = action === 'approve' ? 'approved' : 'rejected';
            this.showMessage('success', `Request ${actionText}`, `Payment request has been ${actionText}`);
            
            // Just remove this specific request from the DOM
            const requestElement = button.closest('.payment-request-item');
            if (requestElement) {
                requestElement.remove();
                
                // Update counter manually
                const incomingCount = document.querySelectorAll('#incomingRequests .payment-request-item').length;
                document.getElementById('incomingCount').textContent = incomingCount.toString();
            }
            
        } catch (error) {
            console.error('Payment request response error:', error);
            this.showMessage('error', 'Action Failed', error.message || `Failed to ${action} payment request`);
        } finally {
            button.innerHTML = originalText;
            button.disabled = false;
        }
    }

    async handlePaymentRequestCancel(requestId, button) {
        const originalText = button.innerHTML;
        button.innerHTML = '<i class="fas fa-spinner fa-spin"></i>';
        button.disabled = true;

        try {
            // Use cancel action for cancelling own requests
            await this.api.handlePaymentRequest(requestId, 'cancel');
            
            this.showMessage('success', 'Request Cancelled', 'Payment request has been cancelled');
            
            // Just remove this specific request from the DOM
            const requestElement = button.closest('.payment-request-item');
            if (requestElement) {
                requestElement.remove();
                
                // Update counter manually
                const outgoingCount = document.querySelectorAll('#outgoingRequests .payment-request-item').length;
                document.getElementById('outgoingCount').textContent = outgoingCount.toString();
            }
            
        } catch (error) {
            console.error('Cancel payment request error:', error);
            this.showMessage('error', 'Cancel Failed', error.message || 'Failed to cancel payment request');
        } finally {
            button.innerHTML = originalText;
            button.disabled = false;
        }
    }

    createTransactionHTML(transaction) {
        const isPositive = transaction.amount > 0;
        const iconClass = 'fa-exchange-alt';
        const iconType = 'transfer';
        
        const formattedDate = this.formatDate(transaction.created_at || transaction.date);
        
        return `
            <div class="transaction-item">
                <div class="transaction-info">
                    <div class="transaction-icon ${iconType}">
                        <i class="fas ${iconClass}"></i>
                    </div>
                    <div class="transaction-details">
                        <h4>${transaction.description}</h4>
                        <p>${formattedDate}</p>
                        ${transaction.memo ? `<p><em>${transaction.memo}</em></p>` : ''}
                    </div>
                </div>
                <div class="transaction-amount ${isPositive ? 'positive' : 'negative'}">
                    ${isPositive ? '+' : ''}$${Math.abs(transaction.amount).toFixed(2)}
                </div>
            </div>
        `;
    }

    openRequestModal() {
        const modal = document.getElementById('requestModal');
        if (modal) {
            modal.classList.add('show');
            document.getElementById('requestForm').reset();
            document.getElementById('requestFromAccount').focus();
        }
    }

    openTransferModal() {
        const modal = document.getElementById('transferModal');
        if (modal) {
            modal.classList.add('show');
            document.getElementById('transferForm').reset();
            document.getElementById('recipientAccount').focus();
        }
    }

    openRegistrationModal() {
        const modal = document.getElementById('registrationModal');
        if (modal) {
            modal.classList.add('show');
            document.getElementById('registrationForm').reset();
            document.getElementById('regUsername').focus();
        }
    }

    openCardModal() {
        const modal = document.getElementById('cardModal');
        if (modal) {
            this.populateCardDetails();
            modal.classList.add('show');
        }
    }

    async populateCardDetails() {
        if (!this.currentUser) return;

        try {
            // Get card data from backend
            const cardResponse = await this.api.getCard();
            const card = cardResponse.card;
            
            // Format card number for display
            const formattedCardNumber = this.formatCardNumber(card.card_number);

            // Update card display
            document.getElementById('cardNumber').textContent = formattedCardNumber;
            document.getElementById('fullCardNumber').textContent = formattedCardNumber;
            
            // Use username since fullName doesn't exist in backend User model
            const cardHolderName = this.currentUser.username.toUpperCase();
            document.getElementById('cardHolderName').textContent = cardHolderName;
            document.getElementById('cardExpiry').textContent = card.expiry_date;

            // Update refresh button state based on backend response
            this.updateRefreshButtonState(cardResponse.canRefresh, cardResponse.timeUntilRefresh);
            
        } catch (error) {
            console.error('Failed to load card details:', error);
            // Fallback to local generation if backend fails
            this.populateCardDetailsLocal();
        }
    }

    populateCardDetailsLocal() {
        if (!this.currentUser) return;

        // Load card refresh state
        this.loadCardRefreshState();

        // Generate account-tied card number
        const cardNumber = this.generateCardNumber(this.currentUser.account_number || this.currentUser.accountNumber);
        const formattedCardNumber = this.formatCardNumber(cardNumber);

        // Generate expiry date (3 years from now)
        const expiryDate = this.generateExpiryDate();

        // Update card display
        document.getElementById('cardNumber').textContent = formattedCardNumber;
        document.getElementById('fullCardNumber').textContent = formattedCardNumber;
        
        // Use username since fullName doesn't exist in backend User model
        const cardHolderName = this.currentUser.username.toUpperCase();
        document.getElementById('cardHolderName').textContent = cardHolderName;
        document.getElementById('cardExpiry').textContent = expiryDate;

        // Update refresh button state
        this.updateRefreshButtonStateLocal();
    }

    generateCardNumber(accountNumber) {
        // Use account number as seed for consistent card generation
        const accountSeed = parseInt(accountNumber.replace(/\D/g, '')) || 1234;
        
        // Add refresh seed to generate new numbers when card is refreshed
        const combinedSeed = accountSeed + this.cardRefreshSeed;
        
        // Viridian City Bank BIN (Bank Identification Number) - using 4532 (Visa format)
        let cardNumber = '4532';
        
        // Generate middle digits based on combined seed
        const middle8Digits = this.generateMiddleDigits(combinedSeed);
        cardNumber += middle8Digits;
        
        // Add check digit using Luhn algorithm
        const checkDigit = this.calculateLuhnCheckDigit(cardNumber);
        cardNumber += checkDigit;
        
        return cardNumber;
    }

    generateMiddleDigits(seed) {
        // Create a pseudo-random generator using the seed
        let random = seed;
        let result = '';
        
        for (let i = 0; i < 8; i++) {
            random = (random * 9301 + 49297) % 233280;
            result += Math.floor((random / 233280) * 10);
        }
        
        return result.padStart(8, '0');
    }

    calculateLuhnCheckDigit(cardNumber) {
        let sum = 0;
        let isEven = true;
        
        // Process digits from right to left
        for (let i = cardNumber.length - 1; i >= 0; i--) {
            let digit = parseInt(cardNumber[i]);
            
            if (isEven) {
                digit *= 2;
                if (digit > 9) {
                    digit = digit - 9;
                }
            }
            
            sum += digit;
            isEven = !isEven;
        }
        
        return (10 - (sum % 10)) % 10;
    }

    formatCardNumber(cardNumber) {
        return cardNumber.replace(/(\d{4})(?=\d)/g, '$1 ');
    }

    generateExpiryDate() {
        const currentDate = new Date();
        const expiryYear = currentDate.getFullYear() + 3;
        const expiryMonth = (currentDate.getMonth() + 1).toString().padStart(2, '0');
        
        return `${expiryMonth}/${expiryYear.toString().slice(-2)}`;
    }

    addToAppleWallet() {
        if (!this.currentUser) return;

        // Check if device supports Apple Wallet
        if (!this.isIOS() && !this.isMac()) {
            this.showMessage('Apple Wallet is only available on iOS and macOS devices.', 'error');
            return;
        }

        // Generate Apple Wallet pass
        const cardData = this.getCardData();
        const passData = this.generateAppleWalletPass(cardData);
        
        // Create download link
        this.downloadWalletPass(passData, 'viridian-bank-card.pkpass');
        this.showMessage('Apple Wallet pass downloaded! Tap to add to your wallet.', 'success');
    }

    addToGooglePay() {
        if (!this.currentUser) return;

        // Check if device supports Google Pay
        if (!this.isAndroid() && !this.isChrome()) {
            this.showMessage('Google Pay integration is best supported on Android devices or Chrome browser.', 'info');
        }

        // Generate Google Pay pass
        const cardData = this.getCardData();
        const passData = this.generateGooglePayPass(cardData);
        
        // Create download link or redirect to Google Pay
        this.downloadWalletPass(passData, 'viridian-bank-card.json');
        this.showMessage('Google Pay pass generated! Follow the instructions to add to your wallet.', 'success');
    }

    getCardData() {
        const cardNumber = this.generateCardNumber(this.currentUser.accountNumber);
        const expiryDate = this.generateExpiryDate();
        
        return {
            cardNumber: cardNumber,
            formattedCardNumber: this.formatCardNumber(cardNumber),
            holderName: this.currentUser.fullName || this.currentUser.username,
            expiryDate: expiryDate,
            bankName: 'Viridian City Bank',
            cardType: 'Debit',
            accountNumber: this.currentUser.accountNumber,
            dailyLimit: '₽50,000.00'
        };
    }

    generateAppleWalletPass(cardData) {
        // Apple Wallet pass structure (simplified for demo)
        const passData = {
            formatVersion: 1,
            passTypeIdentifier: 'pass.com.viridianbank.card',
            serialNumber: cardData.cardNumber,
            teamIdentifier: 'VIRIDIAN',
            organizationName: 'Viridian City Bank',
            description: 'Viridian Bank Debit Card',
            logoText: 'Viridian City Bank',
            foregroundColor: 'rgb(255, 255, 255)',
            backgroundColor: 'rgb(64, 130, 109)',
            generic: {
                primaryFields: [
                    {
                        key: 'balance',
                        label: 'Account Balance',
                        value: document.getElementById('accountBalance').textContent
                    }
                ],
                secondaryFields: [
                    {
                        key: 'cardNumber',
                        label: 'Card Number',
                        value: cardData.formattedCardNumber
                    },
                    {
                        key: 'expires',
                        label: 'Expires',
                        value: cardData.expiryDate
                    }
                ],
                auxiliaryFields: [
                    {
                        key: 'cardHolder',
                        label: 'Card Holder',
                        value: cardData.holderName.toUpperCase()
                    }
                ]
            }
        };

        return JSON.stringify(passData, null, 2);
    }

    generateGooglePayPass(cardData) {
        // Google Pay pass structure (simplified for demo)
        const passData = {
            iss: 'viridian-city-bank',
            aud: 'google',
            typ: 'savetowallet',
            iat: Math.floor(Date.now() / 1000),
            payload: {
                genericObjects: [
                    {
                        id: `${cardData.cardNumber}-${Date.now()}`,
                        classId: 'viridian-bank-card-class',
                        genericType: 'GENERIC_TYPE_UNSPECIFIED',
                        hexBackgroundColor: '#40826D',
                        logo: {
                            sourceUri: {
                                uri: 'https://viridianbank.com/logo.png'
                            }
                        },
                        cardTitle: {
                            defaultValue: {
                                language: 'en',
                                value: 'Viridian Bank Debit Card'
                            }
                        },
                        subheader: {
                            defaultValue: {
                                language: 'en',
                                value: 'Debit Card'
                            }
                        },
                        header: {
                            defaultValue: {
                                language: 'en',
                                value: cardData.holderName.toUpperCase()
                            }
                        },
                        textModulesData: [
                            {
                                id: 'cardNumber',
                                header: 'Card Number',
                                body: cardData.formattedCardNumber
                            },
                            {
                                id: 'expires',
                                header: 'Expires',
                                body: cardData.expiryDate
                            }
                        ]
                    }
                ]
            }
        };

        return JSON.stringify(passData, null, 2);
    }

    downloadWalletPass(passData, filename) {
        const blob = new Blob([passData], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    }

    isIOS() {
        return /iPad|iPhone|iPod/.test(navigator.userAgent);
    }

    isMac() {
        return /Mac|Macintosh/.test(navigator.userAgent);
    }

    isAndroid() {
        return /Android/.test(navigator.userAgent);
    }

    isChrome() {
        return /Chrome/.test(navigator.userAgent);
    }

    loadCardRefreshState() {
        if (!this.currentUser) return;
        
        const storageKey = `cardRefresh_${this.currentUser.username}`;
        const refreshData = localStorage.getItem(storageKey);
        
        if (refreshData) {
            const data = JSON.parse(refreshData);
            this.cardRefreshSeed = data.seed || 0;
            this.lastRefreshDate = data.lastRefresh ? new Date(data.lastRefresh) : null;
        } else {
            this.cardRefreshSeed = 0;
            this.lastRefreshDate = null;
        }
    }

    saveCardRefreshState() {
        if (!this.currentUser) return;
        
        const storageKey = `cardRefresh_${this.currentUser.username}`;
        const refreshData = {
            seed: this.cardRefreshSeed,
            lastRefresh: this.lastRefreshDate ? this.lastRefreshDate.toISOString() : null
        };
        
        localStorage.setItem(storageKey, JSON.stringify(refreshData));
    }

    canRefreshCard() {
        if (!this.lastRefreshDate) return true;
        
        const now = new Date();
        const diffTime = Math.abs(now - this.lastRefreshDate);
        const diffHours = Math.ceil(diffTime / (1000 * 60 * 60));
        
        return diffHours >= 24;
    }

    getTimeUntilNextRefresh() {
        if (!this.lastRefreshDate) return null;
        
        const now = new Date();
        const nextRefreshTime = new Date(this.lastRefreshDate.getTime() + (24 * 60 * 60 * 1000));
        
        if (now >= nextRefreshTime) return null;
        
        const diffTime = nextRefreshTime - now;
        const hours = Math.floor(diffTime / (1000 * 60 * 60));
        const minutes = Math.floor((diffTime % (1000 * 60 * 60)) / (1000 * 60));
        
        return { hours, minutes };
    }

    updateRefreshButtonState(canRefresh, timeUntilRefresh) {
        const refreshBtn = document.getElementById('refreshCardBtn');
        if (!refreshBtn) return;
        
        if (canRefresh) {
            refreshBtn.disabled = false;
            refreshBtn.innerHTML = '<i class="fas fa-sync-alt"></i> Get New Card';
            refreshBtn.title = 'Generate a new card with different numbers';
        } else {
            refreshBtn.disabled = true;
            if (timeUntilRefresh) {
                refreshBtn.innerHTML = `<i class="fas fa-clock"></i> ${timeUntilRefresh.hours}h ${timeUntilRefresh.minutes}m`;
                refreshBtn.title = `You can refresh your card in ${timeUntilRefresh.hours} hours and ${timeUntilRefresh.minutes} minutes`;
            } else {
                refreshBtn.innerHTML = '<i class="fas fa-clock"></i> Please wait';
                refreshBtn.title = 'Card refresh not available at this time';
            }
        }
    }

    updateRefreshButtonStateLocal() {
        const refreshBtn = document.getElementById('refreshCardBtn');
        if (!refreshBtn) return;
        
        const canRefresh = this.canRefreshCard();
        
        if (canRefresh) {
            refreshBtn.disabled = false;
            refreshBtn.innerHTML = '<i class="fas fa-sync-alt"></i> Get New Card';
            refreshBtn.title = 'Generate a new card with different numbers';
        } else {
            refreshBtn.disabled = true;
            const timeUntil = this.getTimeUntilNextRefresh();
            if (timeUntil) {
                refreshBtn.innerHTML = `<i class="fas fa-clock"></i> ${timeUntil.hours}h ${timeUntil.minutes}m`;
                refreshBtn.title = `You can refresh your card in ${timeUntil.hours} hours and ${timeUntil.minutes} minutes`;
            }
        }
    }

    async refreshCard() {
        const refreshBtn = document.getElementById('refreshCardBtn');
        const originalContent = refreshBtn.innerHTML;
        
        // Show loading state
        refreshBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Generating...';
        refreshBtn.disabled = true;

        try {
            // Call backend to refresh card
            const response = await this.api.refreshCard();
            
            // Regenerate card details from backend response
            await this.populateCardDetails();
            
            this.showMessage('success', 'New Card Generated!', 
                'Your card has been refreshed with new numbers and expiry date. The old card is now deactivated.');
                
        } catch (error) {
            console.error('Card refresh error:', error);
            
            // If backend fails, show appropriate error message
            if (error.message.includes('once per day')) {
                this.showMessage('info', 'Card Refresh Limit', error.message);
            } else {
                this.showMessage('error', 'Refresh Failed', 'Unable to generate new card. Please try again.');
            }
            
        } finally {
            refreshBtn.innerHTML = originalContent;
            // Re-populate card details to update button state
            await this.populateCardDetails();
        }
    }

    openPasswordModal() {
        const modal = document.getElementById('passwordModal');
        if (modal) {
            modal.classList.add('show');
            document.getElementById('passwordChangeForm').reset();
            document.getElementById('currentPassword').focus();
        }
    }

    async handleRegistration(e) {
        e.preventDefault();
        
        const username = document.getElementById('regUsername').value;
        const email = document.getElementById('regEmail').value;
        const password = document.getElementById('regPassword').value;
        const confirmPassword = document.getElementById('regConfirmPassword').value;
        const submitBtn = e.target.querySelector('button[type="submit"]');

        // Show loading state
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Creating Account...';
        submitBtn.disabled = true;

        try {
            const response = await this.api.register(username, email, password, confirmPassword);
            
            if (response.success) {
                this.showMessage('success', 'Account Created!', 
                    `Welcome ${response.user.username}! Your account has been created with ${response.user.balance} ₽. You can now log in.`);
                this.closeModal('registrationModal');
                
                // Switch to login form
                document.getElementById('username').value = username;
                document.getElementById('password').focus();
            }
        } catch (error) {
            this.showMessage('error', 'Registration Failed', error.message);
        } finally {
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    }

    async handlePasswordChange(e) {
        e.preventDefault();
        
        const currentPassword = document.getElementById('currentPassword').value;
        const newPassword = document.getElementById('newPassword').value;
        const confirmNewPassword = document.getElementById('confirmNewPassword').value;
        const submitBtn = e.target.querySelector('button[type="submit"]');

        // Show loading state
        const originalText = submitBtn.innerHTML;
        submitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Changing Password...';
        submitBtn.disabled = true;

        try {
            const response = await this.api.changePassword(currentPassword, newPassword, confirmNewPassword);
            
            if (response.success) {
                this.showMessage('success', 'Password Changed!', response.message);
                this.closeModal('passwordModal');
                
                // Log out the user
                setTimeout(() => {
                    this.handleLogout();
                }, 2000);
            }
        } catch (error) {
            this.showMessage('error', 'Password Change Failed', error.message);
        } finally {
            submitBtn.innerHTML = originalText;
            submitBtn.disabled = false;
        }
    }

    closeModal(modalId) {
        const modal = document.getElementById(modalId);
        if (modal) {
            modal.classList.remove('show');
        }
    }

    showMessage(type, title, message) {
        const container = document.getElementById('messageContainer');
        if (!container) return;

        const messageEl = document.createElement('div');
        messageEl.className = `message ${type}`;
        messageEl.innerHTML = `
            <h4>${title}</h4>
            <p>${message}</p>
        `;

        container.appendChild(messageEl);

        // Auto remove after 5 seconds
        setTimeout(() => {
            if (messageEl.parentNode) {
                messageEl.style.transform = 'translateX(100%)';
                messageEl.style.opacity = '0';
                setTimeout(() => {
                    if (messageEl.parentNode) {
                        messageEl.parentNode.removeChild(messageEl);
                    }
                }, 300);
            }
        }, 5000);
    }
}

// Initialize the app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.bankingApp = new BankingApp();
});

// Handle browser back/forward buttons
window.addEventListener('popstate', (event) => {
    // Handle navigation if needed
});
