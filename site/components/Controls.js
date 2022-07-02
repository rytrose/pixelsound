import Modes from "./modes/Modes";

const Controls = ({ onModeChange }) => {
  return (
    <div className="py-3">
      <Modes onChange={onModeChange}></Modes>
    </div>
  );
};

export default Controls;
