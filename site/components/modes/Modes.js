const RadioOption = ({
  name,
  value,
  children,
  onChange,
  defaultChecked = false,
}) => {
  const id = `radio-option-${value}`;

  return (
    <div className="flex grow basis-0 justify-center">
      <input
        type="radio"
        name={name}
        id={id}
        value={value}
        onChange={onChange}
        className="peer opacity-0 absolute left-[-99999rem]"
        defaultChecked={defaultChecked}
      ></input>
      <label
        htmlFor={id}
        className="outline-none select-none cursor-pointer 
       p-2 rounded-xl text-black
      bg-violet-300 hover:bg-violet-400
      peer-checked:bg-violet-600 active:bg-violet-600"
      >
        {children}
      </label>
    </div>
  );
};

const Modes = ({ onChange }) => {
  return (
    <div>
      <h2
        className="flex justify-center pb-2
        font-serif text-xl"
      >
        Mode
      </h2>
      <div className="flex sm:max-w-sm sm:mx-auto">
        <RadioOption
          name="mode"
          value="mouse"
          onChange={onChange}
          defaultChecked={true}
        >
          Mouse
        </RadioOption>
        <RadioOption name="mode" value="keyboard" onChange={onChange}>
          Keyboard
        </RadioOption>
        <RadioOption name="mode" value="algorithm" onChange={onChange}>
          Algorithm
        </RadioOption>
      </div>
    </div>
  );
};

export default Modes;
