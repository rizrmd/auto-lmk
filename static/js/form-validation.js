/**
 * Form Validation Framework for Auto LMK Admin Interface
 * Provides consistent validation across all admin forms
 */

class FormValidator {
    constructor(formId, options = {}) {
        this.form = document.getElementById(formId);
        this.options = {
            validateOnBlur: true,
            validateOnInput: false,
            showErrorIcons: true,
            language: 'id',
            ...options
        };
        this.errors = {};
        this.validators = {};
        this.init();
    }

    init() {
        if (!this.form) {
            console.error(`Form with ID "${formId}" not found`);
            return;
        }

        this.setupValidators();
        this.attachEventListeners();
    }

    setupValidators() {
        // Common validators
        this.validators.phone = {
            pattern: /^(\+?[1-9]\d{1,14}|08\d{8,11})$/,
            message: 'Format nomor telepon tidak valid. Gunakan +628xxx atau 08xxx',
            required: true
        };

        this.validators.email = {
            pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
            message: 'Format email tidak valid',
            required: false
        };

        this.validators.name = {
            pattern: /^[a-zA-Z\s]+$/,
            minLength: 2,
            maxLength: 100,
            message: 'Nama hanya boleh berisi huruf dan spasi',
            required: true
        };

        this.validators.required = {
            message: 'Field ini wajib diisi',
            required: true
        };

        this.validators.minLength = (length) => ({
            message: `Minimal ${length} karakter`,
            validate: (value) => value.length >= length
        });

        this.validators.maxLength = (length) => ({
            message: `Maksimal ${length} karakter`,
            validate: (value) => value.length <= length
        });
    }

    attachEventListeners() {
        const inputs = this.form.querySelectorAll('input, textarea, select');

        inputs.forEach(input => {
            // Validate on blur
            if (this.options.validateOnBlur) {
                input.addEventListener('blur', () => {
                    this.validateField(input);
                });
            }

            // Validate on input (for specific fields)
            if (this.options.validateOnInput) {
                input.addEventListener('input', () => {
                    if (input.dataset.validateRealtime === 'true') {
                        this.validateField(input);
                    }
                });
            }

            // Clear error on input
            input.addEventListener('input', () => {
                this.clearFieldError(input.name);
            });
        });

        // Validate on form submit
        this.form.addEventListener('submit', (e) => {
            if (!this.validateForm()) {
                e.preventDefault();
                return false;
            }
        });
    }

    validateField(input) {
        const fieldName = input.name;
        const value = input.value.trim();
        let isValid = true;
        let errorMessage = '';

        // Check required
        if (input.required && !value) {
            isValid = false;
            errorMessage = `${this.getFieldLabel(fieldName)} wajib diisi`;
        }

        // Check pattern
        if (isValid && value && input.pattern) {
            const pattern = new RegExp(input.pattern);
            if (!pattern.test(value)) {
                isValid = false;
                errorMessage = this.getPatternMessage(fieldName) || 'Format tidak valid';
            }
        }

        // Check custom validators from data attributes
        if (isValid && value) {
            const validators = input.dataset.validators ? input.dataset.validators.split(',') : [];
            validators.forEach(validatorName => {
                const validator = this.getValidator(validatorName);
                if (validator && validator.validate && !validator.validate(value)) {
                    isValid = false;
                    errorMessage = validator.message;
                }
            });
        }

        // Check min/max length
        if (isValid && value) {
            if (input.minLength && value.length < parseInt(input.minLength)) {
                isValid = false;
                errorMessage = `Minimal ${input.minLength} karakter`;
            }
            if (input.maxLength && value.length > parseInt(input.maxLength)) {
                isValid = false;
                errorMessage = `Maksimal ${input.maxLength} karakter`;
            }
        }

        if (!isValid) {
            this.showFieldError(fieldName, errorMessage);
        } else {
            this.clearFieldError(fieldName);
        }

        return isValid;
    }

    validateForm() {
        const inputs = this.form.querySelectorAll('input, textarea, select');
        let isValid = true;

        this.errors = {};

        inputs.forEach(input => {
            if (!this.validateField(input)) {
                isValid = false;
            }
        });

        // Show form-level error if any
        if (!isValid) {
            this.showFormError('Silakan perbaiki error yang ditandai');
        } else {
            this.clearFormError();
        }

        return isValid;
    }

    showFieldError(fieldName, message) {
        this.errors[fieldName] = message;
        const input = this.form.querySelector(`[name="${fieldName}"]`);

        if (input) {
            // Add error class to input
            input.classList.add('border-red-500', 'focus:ring-red-500');
            input.classList.remove('border-gray-300', 'focus:ring-blue-500');

            // Create or update error message element
            let errorElement = input.parentNode.querySelector('.field-error');
            if (!errorElement) {
                errorElement = document.createElement('p');
                errorElement.className = 'field-error text-red-600 text-sm mt-1 flex items-center';
                input.parentNode.appendChild(errorElement);
            }

            errorElement.innerHTML = this.options.showErrorIcons
                ? `<span class="mr-1">⚠️</span>${message}`
                : message;
        }
    }

    clearFieldError(fieldName) {
        delete this.errors[fieldName];
        const input = this.form.querySelector(`[name="${fieldName}"]`);

        if (input) {
            // Remove error class from input
            input.classList.remove('border-red-500', 'focus:ring-red-500');
            input.classList.add('border-gray-300', 'focus:ring-blue-500');

            // Remove error message element
            const errorElement = input.parentNode.querySelector('.field-error');
            if (errorElement) {
                errorElement.remove();
            }
        }
    }

    showFormError(message) {
        let errorElement = this.form.querySelector('.form-error');
        if (!errorElement) {
            errorElement = document.createElement('div');
            errorElement.className = 'form-error mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded';
            this.form.insertBefore(errorElement, this.form.firstChild);
        }
        errorElement.textContent = message;
    }

    clearFormError() {
        const errorElement = this.form.querySelector('.form-error');
        if (errorElement) {
            errorElement.remove();
        }
    }

    getValidator(name) {
        return this.validators[name];
    }

    getFieldLabel(fieldName) {
        const label = this.form.querySelector(`label[for="${fieldName}"]`);
        if (label) {
            return label.textContent.replace('*', '').trim();
        }
        return fieldName.charAt(0).toUpperCase() + fieldName.slice(1);
    }

    getPatternMessage(fieldName) {
        const input = this.form.querySelector(`[name="${fieldName}"]`);
        if (input && input.dataset.validationMessage) {
            return input.dataset.validationMessage;
        }
        return null;
    }

    // Static method to auto-initialize forms with data-form-validator attribute
    static autoInit() {
        document.querySelectorAll('[data-form-validator]').forEach(form => {
            const formId = form.id || `form-${Date.now()}`;
            if (!form.id) form.id = formId;

            const options = {};
            if (form.dataset.validateOnBlur === 'false') options.validateOnBlur = false;
            if (form.dataset.validateOnInput === 'true') options.validateOnInput = true;

            new FormValidator(formId, options);
        });
    }
}

// Auto-initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    FormValidator.autoInit();
});

// Make available globally
window.FormValidator = FormValidator;