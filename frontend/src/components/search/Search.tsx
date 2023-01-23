import React from "react";

type SearchProps = {
  value?: string;
  placeholder?: string;
  disabled?: boolean;
  autocomplete?: boolean;
  onChange?: React.ChangeEventHandler<HTMLInputElement>;
};

export const OnlyofficeSearchBar: React.FC<SearchProps> = ({
  value,
  placeholder,
  disabled,
  autocomplete = false,
  onChange,
}) => (
  <div className="font-sans text-black bg-white w-screen">
    <div className="border rounded overflow-hidden flex">
      <input
        type="text"
        className="py-2 px-2 w-full select-auto outline-none"
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        disabled={disabled}
        autoCorrect={autocomplete ? undefined : "off"}
        autoComplete={autocomplete ? undefined : "off"}
      />
      <button
        type="button"
        className={`px-6 ${disabled && "bg-gray-50"}`}
        disabled={disabled}
      >
        <svg
          className="h-4 w-4 text-grey-dark"
          fill="currentColor"
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
        >
          <path d="M16.32 14.9l5.39 5.4a1 1 0 0 1-1.42 1.4l-5.38-5.38a8 8 0 1 1 1.41-1.41zM10 16a6 6 0 1 0 0-12 6 6 0 0 0 0 12z" />
        </svg>
      </button>
    </div>
  </div>
);
