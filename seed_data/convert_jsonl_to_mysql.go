package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Occupation struct {
	ID              string      `json:"id"`
	SocID           string      `json:"socId"`
	SocTitle        string      `json:"socTitle"`
	Title           string      `json:"title"`
	SingularTitle   string      `json:"singularTitle"`
	Description     string      `json:"description"`
	TypicalEdLevel  string      `json:"typicalEdLevel"`
	CoreTasks       []string    `json:"coreTasks"`
	Skills          []Skill     `json:"skills"`
	Knowledge       []Knowledge `json:"knowledge"`
	Abilities       []Ability   `json:"abilities"`
}

type Skill struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Importance  float64 `json:"importance"`
	Level       float64 `json:"level"`
}

type Knowledge struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Importance  float64 `json:"importance"`
	Level       float64 `json:"level"`
}

type Ability struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Importance  float64 `json:"importance"`
	Level       float64 `json:"level"`
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "''")
	return "'" + s + "'"
}

func writeSchema(f *os.File) {
	schema := `-- Database schema for occupations data
CREATE TABLE IF NOT EXISTS occupations (
    id VARCHAR(20) PRIMARY KEY,
    soc_id VARCHAR(20),
    soc_title VARCHAR(255),
    title VARCHAR(255),
    singular_title VARCHAR(255),
    description TEXT,
    typical_ed_level VARCHAR(100),
    data JSON,
    INDEX idx_soc_id (soc_id),
    INDEX idx_title (title)
);

CREATE TABLE IF NOT EXISTS occupation_tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    occupation_id VARCHAR(20),
    task TEXT,
    FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE CASCADE,
    INDEX idx_occupation_id (occupation_id)
);

CREATE TABLE IF NOT EXISTS occupation_skills (
    id INT AUTO_INCREMENT PRIMARY KEY,
    occupation_id VARCHAR(20),
    skill_name VARCHAR(255),
    skill_description TEXT,
    importance DECIMAL(3,2),
    level DECIMAL(10,6),
    FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE CASCADE,
    INDEX idx_occupation_id (occupation_id),
    INDEX idx_skill_name (skill_name)
);

CREATE TABLE IF NOT EXISTS occupation_knowledge (
    id INT AUTO_INCREMENT PRIMARY KEY,
    occupation_id VARCHAR(20),
    knowledge_name VARCHAR(255),
    knowledge_description TEXT,
    importance DECIMAL(3,2),
    level DECIMAL(10,6),
    FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE CASCADE,
    INDEX idx_occupation_id (occupation_id),
    INDEX idx_knowledge_name (knowledge_name)
);

CREATE TABLE IF NOT EXISTS occupation_abilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    occupation_id VARCHAR(20),
    ability_name VARCHAR(255),
    ability_description TEXT,
    importance DECIMAL(3,2),
    level DECIMAL(10,6),
    FOREIGN KEY (occupation_id) REFERENCES occupations(id) ON DELETE CASCADE,
    INDEX idx_occupation_id (occupation_id),
    INDEX idx_ability_name (ability_name)
);

`
	f.WriteString(schema)
}

func generateSQL(occ Occupation, rawJSON []byte) string {
	var sql strings.Builder

	// Main occupation insert
	sql.WriteString(fmt.Sprintf(
		"INSERT INTO occupations (id, soc_id, soc_title, title, singular_title, description, typical_ed_level, data)\nVALUES (%s, %s, %s, %s, %s, %s, %s, %s);\n\n",
		escapeString(occ.ID),
		escapeString(occ.SocID),
		escapeString(occ.SocTitle),
		escapeString(occ.Title),
		escapeString(occ.SingularTitle),
		escapeString(occ.Description),
		escapeString(occ.TypicalEdLevel),
		escapeString(string(rawJSON)),
	))

	// Insert tasks
	for _, task := range occ.CoreTasks {
		sql.WriteString(fmt.Sprintf(
			"INSERT INTO occupation_tasks (occupation_id, task) VALUES (%s, %s);\n",
			escapeString(occ.ID),
			escapeString(task),
		))
	}
	if len(occ.CoreTasks) > 0 {
		sql.WriteString("\n")
	}

	// Insert skills
	for _, skill := range occ.Skills {
		sql.WriteString(fmt.Sprintf(
			"INSERT INTO occupation_skills (occupation_id, skill_name, skill_description, importance, level) VALUES (%s, %s, %s, %.2f, %.6f);\n",
			escapeString(occ.ID),
			escapeString(skill.Name),
			escapeString(skill.Description),
			skill.Importance,
			skill.Level,
		))
	}
	if len(occ.Skills) > 0 {
		sql.WriteString("\n")
	}

	// Insert knowledge
	for _, know := range occ.Knowledge {
		sql.WriteString(fmt.Sprintf(
			"INSERT INTO occupation_knowledge (occupation_id, knowledge_name, knowledge_description, importance, level) VALUES (%s, %s, %s, %.2f, %.6f);\n",
			escapeString(occ.ID),
			escapeString(know.Name),
			escapeString(know.Description),
			know.Importance,
			know.Level,
		))
	}
	if len(occ.Knowledge) > 0 {
		sql.WriteString("\n")
	}

	// Insert abilities
	for _, ability := range occ.Abilities {
		sql.WriteString(fmt.Sprintf(
			"INSERT INTO occupation_abilities (occupation_id, ability_name, ability_description, importance, level) VALUES (%s, %s, %s, %.2f, %.6f);\n",
			escapeString(occ.ID),
			escapeString(ability.Name),
			escapeString(ability.Description),
			ability.Importance,
			ability.Level,
		))
	}
	if len(occ.Abilities) > 0 {
		sql.WriteString("\n")
	}

	return sql.String()
}

func convertJSONLToSQL(inputFile, outputFile string, maxLines int) error {
	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer input.Close()

	output, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer output.Close()

	// Write schema
	writeSchema(output)
	output.WriteString("-- Data inserts\n\n")

	scanner := bufio.NewScanner(input)
	lineCount := 0

	for scanner.Scan() {
		if maxLines > 0 && lineCount >= maxLines {
			break
		}

		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var occ Occupation
		if err := json.Unmarshal(line, &occ); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid JSON on line %d: %v\n", lineCount+1, err)
			continue
		}

		sql := generateSQL(occ, line)
		output.WriteString(sql)
		output.WriteString("-- ------------------------------------------------\n\n")

		lineCount++
		if lineCount%10 == 0 {
			fmt.Printf("Processed %d records...\n", lineCount)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	fmt.Printf("\nConversion complete! Processed %d records.\n", lineCount)
	fmt.Printf("Output written to: %s\n", outputFile)

	return nil
}

func main() {
	inputFile := flag.String("input", "occupations.jsonl", "Input JSONL file")
	outputFile := flag.String("output", "seed_data.sql", "Output SQL file")
	maxLines := flag.Int("lines", 100, "Maximum number of lines to process (0 for all)")

	flag.Parse()

	if err := convertJSONLToSQL(*inputFile, *outputFile, *maxLines); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}