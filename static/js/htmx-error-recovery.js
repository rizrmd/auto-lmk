/**
 * HTMX Error Recovery and Loading States Framework
 * Provides consistent loading indicators and error recovery for HTMX requests
 */

class HtmxErrorRecovery {
    constructor() {
        this.init();
    }

    init() {
        this.setupGlobalLoading();
        this.setupErrorHandling();
        this.setupLoadingIndicators();
        this.setupRetryMechanisms();
    }

    setupGlobalLoading() {
        // Global loading bar at top of page
        if (!document.getElementById('global-loading')) {
            const loadingBar = document.createElement('div');
            loadingBar.id = 'global-loading';
            loadingBar.className = 'fixed top-0 left-0 w-full h-1 bg-blue-500 transform -translate-x-full transition-transform duration-300 z-50';
            loadingBar.innerHTML = `
                <div class="h-full bg-blue-600 animate-pulse"></div>
            `;
            document.body.appendChild(loadingBar);
        }
    }

    setupErrorHandling() {
        // Handle HTMX errors globally
        document.body.addEventListener('htmx:responseError', function(evt) {
            const xhr = evt.detail.xhr;
            const target = evt.detail.target;

            console.error('HTMX Error:', {
                status: xhr.status,
                response: xhr.responseText,
                target: target
            });

            // Remove loading states
            this.hideLoading();

            // Handle different error types
            if (xhr.status === 0) {
                // Network error
                this.showNetworkError(target, 'Tidak dapat terhubung ke server. Silakan periksa koneksi internet Anda.');
            } else if (xhr.status >= 500) {
                // Server error
                this.showServerError(target, xhr.responseText);
            } else if (xhr.status >= 400) {
                // Client error
                this.showClientError(target, xhr.responseText);
            } else {
                // Unknown error
                this.showGenericError(target, 'Terjadi kesalahan yang tidak diketahui');
            }
        });

        // Handle HTMX timeouts
        document.body.addEventListener('htmx:timeout', function(evt) {
            const target = evt.detail.target;
            this.hideLoading();
            this.showTimeoutError(target, 'Request terlalu lama. Silakan coba lagi.');
        }.bind(this));

        // Handle successful requests
        document.body.addEventListener('htmx:afterRequest', function(evt) {
            const xhr = evt.detail.xhr;

            if (xhr.status < 400) {
                // Success - any additional success handling can go here
                this.hideLoading();
            }
        }.bind(this));

        // Show loading on request
        document.body.addEventListener('htmx:beforeRequest', function(evt) {
            this.showLoading();
        }.bind(this));
    }

    setupLoadingIndicators() {
        // Add loading class to buttons during HTMX requests
        document.body.addEventListener('htmx:beforeRequest', function(evt) {
            const target = evt.detail.target;

            // Find the button that triggered the request
            const button = target.closest('button') || target.querySelector('button');
            if (button) {
                this.setButtonLoading(button, true);
            }

            // Find form inputs
            const form = target.closest('form');
            if (form) {
                this.setFormLoading(form, true);
            }
        }.bind(this));

        // Remove loading class after request
        document.body.addEventListener('htmx:afterRequest', function(evt) {
            const target = evt.detail.target;

            // Find the button that triggered the request
            const button = target.closest('button') || target.querySelector('button');
            if (button) {
                this.setButtonLoading(button, false);
            }

            // Find form inputs
            const form = target.closest('form');
            if (form) {
                this.setFormLoading(form, false);
            }
        }.bind(this));
    }

    setupRetryMechanisms() {
        // Add retry functionality to error messages
        document.body.addEventListener('click', function(evt) {
            if (evt.target.matches('.htmx-retry-btn')) {
                const originalRequest = evt.target.dataset.originalRequest;
                if (originalRequest) {
                    this.retryRequest(originalRequest);
                }
            }
        }.bind(this));
    }

    showLoading() {
        const loadingBar = document.getElementById('global-loading');
        if (loadingBar) {
            loadingBar.classList.remove('-translate-x-full');
        }
    }

    hideLoading() {
        const loadingBar = document.getElementById('global-loading');
        if (loadingBar) {
            loadingBar.classList.add('-translate-x-full');
        }
    }

    setButtonLoading(button, loading) {
        if (loading) {
            button.classList.add('opacity-50', 'cursor-not-allowed', 'disabled');
            button.disabled = true;

            // Store original content
            const originalContent = button.innerHTML;
            button.dataset.originalContent = originalContent;

            // Show loading state
            button.innerHTML = `
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Loading...
            `;
        } else {
            button.classList.remove('opacity-50', 'cursor-not-allowed', 'disabled');
            button.disabled = false;

            // Restore original content
            if (button.dataset.originalContent) {
                button.innerHTML = button.dataset.originalContent;
                delete button.dataset.originalContent;
            }
        }
    }

    setFormLoading(form, loading) {
        const inputs = form.querySelectorAll('input, textarea, select, button');
        inputs.forEach(input => {
            if (loading) {
                input.classList.add('opacity-50');
                if (input.tagName === 'BUTTON' || input.type === 'submit') {
                    input.disabled = true;
                }
            } else {
                input.classList.remove('opacity-50');
                input.disabled = false;
            }
        });
    }

    showNetworkError(target, message) {
        this.showError(target, message, 'error', true);
    }

    showServerError(target, response) {
        let message = 'Server sedang bermasalah. Silakan coba lagi nanti.';

        try {
            const errorData = JSON.parse(response);
            if (errorData.error) {
                message = errorData.error;
            }
        } catch (e) {
            // Use default message if JSON parsing fails
        }

        this.showError(target, message, 'error', true);
    }

    showClientError(target, response) {
        let message = 'Terjadi kesalahan pada permintaan Anda.';

        try {
            const errorData = JSON.parse(response);
            if (errorData.error) {
                message = errorData.error;
            }
        } catch (e) {
            // Use default message if JSON parsing fails
        }

        this.showError(target, message, 'warning');
    }

    showTimeoutError(target, message) {
        this.showError(target, message, 'error', true);
    }

    showGenericError(target, message) {
        this.showError(target, message, 'error');
    }

    showError(target, message, type = 'error', showRetry = false) {
        // Remove existing error messages
        this.removeErrorMessages(target);

        // Create error element
        const errorDiv = document.createElement('div');
        errorDiv.className = `htmx-error-message p-3 mb-4 rounded border ${
            type === 'error' ? 'bg-red-100 border-red-400 text-red-700' :
            type === 'warning' ? 'bg-yellow-100 border-yellow-400 text-yellow-700' :
            'bg-blue-100 border-blue-400 text-blue-700'
        }`;

        let errorContent = `
            <div class="flex items-center">
                <span class="mr-2">
                    ${type === 'error' ? '❌' : type === 'warning' ? '⚠️' : 'ℹ️'}
                </span>
                <span>${message}</span>
            </div>
        `;

        if (showRetry) {
            errorContent += `
                <div class="mt-2">
                    <button class="htmx-retry-btn text-sm underline hover:no-underline" data-original-request="${this.getLastRequestData()}">
                        Coba Lagi
                    </button>
                </div>
            `;
        }

        errorDiv.innerHTML = errorContent;

        // Insert error message
        if (target) {
            target.insertBefore(errorDiv, target.firstChild);
        } else {
            // Fallback to body
            this.showToast(message, type);
        }

        // Auto-remove after 10 seconds
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.remove();
            }
        }, 10000);
    }

    removeErrorMessages(target) {
        if (target) {
            const errors = target.querySelectorAll('.htmx-error-message');
            errors.forEach(error => error.remove());
        }
    }

    showToast(message, type = 'info') {
        // Remove existing toasts
        const existingToast = document.querySelector('.htmx-toast');
        if (existingToast) {
            existingToast.remove();
        }

        // Create toast
        const toast = document.createElement('div');
        toast.className = `htmx-toast fixed top-4 right-4 px-6 py-3 rounded shadow-lg text-white z-50 transition-opacity duration-300 ${
            type === 'error' ? 'bg-red-500' :
            type === 'warning' ? 'bg-yellow-500' :
            type === 'success' ? 'bg-green-500' :
            'bg-blue-500'
        }`;
        toast.textContent = message;

        document.body.appendChild(toast);

        // Auto-remove after 3 seconds
        setTimeout(() => {
            toast.style.opacity = '0';
            setTimeout(() => {
                if (toast.parentNode) {
                    toast.parentNode.removeChild(toast);
                }
            }, 300);
        }, 3000);
    }

    getLastRequestData() {
        // Store last request data for retry functionality
        // This is a simplified version - in production, you'd want more sophisticated tracking
        return localStorage.getItem('lastHtmxRequest') || '';
    }

    retryRequest(requestData) {
        if (requestData) {
            // Implement retry logic
            console.log('Retrying request:', requestData);
            // This would need more sophisticated implementation based on request type
            this.showToast('Mencoba ulang permintaan...', 'info');
        }
    }
}

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    window.htmxErrorRecovery = new HtmxErrorRecovery();
});

// Make available globally
window.HtmxErrorRecovery = HtmxErrorRecovery;