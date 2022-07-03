import FileInput from "./FileInput";
import Modes from "./modes/Modes";

const Controls = ({ onImageChange, onAudioChange, onModeChange }) => {
  return (
    <div className="flex flex-col py-3">
      <div className="flex gap-4 pl-3 pb-3">
        <FileInput onChange={onImageChange} accept=".jpg,.jpeg,.png">
          Select an image
        </FileInput>
        <FileInput onChange={onAudioChange} accept=".mp3">
          Select an audio file
        </FileInput>
      </div>
      <Modes onChange={onModeChange}></Modes>
    </div>
  );
};

export default Controls;
