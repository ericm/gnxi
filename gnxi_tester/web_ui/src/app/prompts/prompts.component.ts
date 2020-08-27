import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { FileService } from '../file.service';
import { Prompts, PromptsList, PromptsSet } from '../models/Prompts';
import { PromptsService } from '../prompts.service';

type ControlGroup = { [name: string]: FormControl }

@Component({
  selector: 'app-prompts',
  templateUrl: './prompts.component.html',
  styleUrls: ['./prompts.component.css']
})
export class PromptsComponent implements OnInit {
  prompts: PromptsSet = {};
  promptsList: PromptsList = { prompts: [], files: [] };
  controlGroup = this.formBuilder.group({ name: ['', Validators.required] });
  files = {};

  constructor(public promptsService: PromptsService, private formBuilder: FormBuilder, private snackbar: MatSnackBar, private fileService: FileService) {
    this.init();
  }

  ngOnInit(): void {

  }

  async init() {
    try {
      this.promptsList = await this.promptsService.getPromptsList().toPromise();
      this.prompts = await this.promptsService.getPrompts().toPromise();
      let fields = this.getFields();
      this.controlGroup = this.formBuilder.group(fields);
    } catch (e) {
      console.error(e);
    }
  }

  private getFields(): { [name: string]: any } {
    let fields = { name: ['', Validators.required] };
    for (let field of this.promptsList.prompts) {
      fields["prompts_" + field] = ['', Validators.required];
    }
    for (let field of this.promptsList.files) {
      fields["files_" + field] = [''];
    }
    return fields;
  }

  setFile(name: string, val: string) {
    let fields = {};
    fields[`files_${name}`] = val;
    this.controlGroup.patchValue(fields)
  }

  deletePrompts(): void {
    let name = this.controlGroup.get("name").value;
    this.promptsService.delete(name).subscribe(res => {
      console.log(res);
      delete this.prompts[name];
      this.controlGroup.reset()
      this.snackbar.open("Deleted", "", {duration: 2000});
    }, error => console.error(error))
  }

  setPrompts(form: { [key: string]: string }): void {
    if (!form) {
      return;
    }
    console.log(form);
    let prompts: Prompts = {
      name: form.name,
      prompts: {},
      files: {}
    }
    for (let field of Object.keys(form)) {
      if (field.search("prompts_") === 0) {
        let key = field.slice(8);
        prompts.prompts[key] = form[field];
      } else if (field.search("files_") === 0) {
        let key = field.slice(6);
        prompts.files[key] = form[field];
      }
    }
    this.promptsService.setPrompts(prompts).subscribe(
      (res) => {
        console.log(res);
        this.snackbar.open("Saved", "", { duration: 2000 });
        this.setSelectedPrompts(prompts.name);
      },
      (err) => console.log(err)
    );
  }

  setSelectedPrompts(name: string): void {
    const prompts = this.prompts[name];
    if (prompts === undefined) {
      this.controlGroup.reset();
      return;
    }
    let fields = {};
    console.log(prompts.files)
    for (let field of Object.keys(prompts.files)) {
      fields[`files_${field}`] = prompts.files[field];
    }
    for (let field of Object.keys(prompts.prompts)) {
      fields[`prompts_${field}`] = prompts.prompts[field];
    }
    this.controlGroup.setValue({
      name,
      ...fields,
    });
  }

  get selected() {
    return this.controlGroup.get("name").value
  }
}
